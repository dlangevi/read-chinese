package core

import (
	"database/sql"
	"embed"
	"fmt"
	"log"
	"os"
  "errors"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/sqlite3"
	"github.com/golang-migrate/migrate/v4/source/iofs"
	_ "github.com/mattn/go-sqlite3"
)

//go:embed migrations
var fs embed.FS

func NewDB(dbPath string) (*sql.DB, error) {
  if _, err := os.Stat(dbPath); err == nil {
    // path/to/whatever exists

} else if errors.Is(err, os.ErrNotExist) {
    // path/to/whatever does *not* exist
	  // os.MkdirAll("./data/1234", 0755)
    os.Create(dbPath)

} else {
    // Schrodinger: file may or may not exist. See err for details.
    // Therefore, do *NOT* use !os.IsNotExist(err) to test for file existence
    log.Fatal(err);
  }

	sqliteDb, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		log.Fatal(err, "failed to open sqlite DB")
	}

	return sqliteDb, nil
}

func RunMigrateScripts(db *sql.DB) error {
	driver, err := sqlite3.WithInstance(db, &sqlite3.Config{})
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
