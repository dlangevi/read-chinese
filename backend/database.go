package backend

import (
	"embed"
	"errors"
	"fmt"
	"log"
	"os"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/sqlite3"
	"github.com/golang-migrate/migrate/v4/source/iofs"
	"github.com/jmoiron/sqlx"
	_ "github.com/mattn/go-sqlite3"
)

//go:embed migrations
var fs embed.FS

func NewDB(dbPath string) (*sqlx.DB, error) {
	if _, err := os.Stat(dbPath); err == nil {
		// path/to/whatever exists

	} else if errors.Is(err, os.ErrNotExist) {
		// path/to/whatever does *not* exist
		os.Create(dbPath)

	} else {
		// Schrodinger: file may or may not exist. See err for details.
		// Therefore, do *NOT* use !os.IsNotExist(err) to test for file existence
		log.Fatal(err)
	}

	// sqliteDb, err := sql.Open("sqlite3", dbPath)
	sqliteDb, err := sqlx.Connect("sqlite3", dbPath)
	if err != nil {
		return nil, err
	}

	return sqliteDb, nil
}

func RunMigrateScripts(db *sqlx.DB) error {
	driver, err := sqlite3.WithInstance(db.DB, &sqlite3.Config{})
	if err != nil {
		return fmt.Errorf("creating sqlite3 db driver failed %s", err)
	}

	d, err := iofs.New(fs, "migrations")
	if err != nil {
		log.Fatal(err)
	}

	m, err := migrate.NewWithInstance("iofs", d, "sqlite3", driver)
	if err != nil {
		return fmt.Errorf("initializing db migration failed %s", err)
	}

	err = m.Up()
	if err != nil && err != migrate.ErrNoChange {
		return fmt.Errorf("migrating database failed %s", err)
	}

	return nil
}
