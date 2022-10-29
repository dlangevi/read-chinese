package core 

import (
	"database/sql"
	// "fmt"
  "log"

	// "github.com/golang-migrate/migrate"
	// "github.com/golang-migrate/migrate/database/sqlite3"
	_ "github.com/mattn/go-sqlite3"
)

func NewDB(dbPath string) (*sql.DB, error) {
	sqliteDb, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		log.Fatal(err, "failed to open sqlite DB")
	}

	return sqliteDb, nil
}

// func RunMigrateScripts(db *sql.DB) error {
// 	driver, err := sqlite3.WithInstance(db, &sqlite3.Config{})
// 	if err != nil {
// 		return fmt.Errorf("creating sqlite3 db driver failed %s", err)
// 	}
//
// 	res := bindata.Resource(AssetNames(),
// 		func(name string) ([]byte, error) {
// 			return Asset(name)
// 		})
//
// 	d, err := bindata.WithInstance(res)
// 	m, err := migrate.NewWithInstance("go-bindata", d, "sqlite3", driver)
// 	if err != nil {
// 		return fmt.Errorf("initializing db migration failed %s", err)
// 	}
//
// 	err = m.Up()
// 	if err != nil && err != migrate.ErrNoChange {
// 		return fmt.Errorf("migrating database failed %s", err)
// 	}
//
// 	return nil
// }
//
