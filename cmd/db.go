package cmd

// Database actions here

import (
	"database/sql"
	"night/cmd/flags"

	_ "github.com/lib/pq"
	_ "github.com/mattn/go-sqlite3"
)

type DataBase struct {
	DB               *sql.DB
	DBDriver         *flags.DataBaseDriver
	SqlFilePath      string
	ConnectionString string
}

func OpenDB(driver flags.DataBaseDriver, SqlFilePath, ConnectionString string) (*DataBase, error) {
	db, err := sql.Open(driver.String(), ConnectionString)
	if err != nil {
		return nil, err
	}

	newDb := DataBase{
		DB:               db,
		DBDriver:         &driver,
		SqlFilePath:      SqlFilePath,
		ConnectionString: ConnectionString,
	}

	return &newDb, nil
}
