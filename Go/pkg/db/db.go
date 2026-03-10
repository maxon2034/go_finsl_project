package db

import (
	"database/sql"

	_ "modernc.org/sqlite"
)

const schema = "CREATE TABLE IF NOT EXISTS scheduler (" +
	"id INTEGER PRIMARY KEY AUTOINCREMENT," +
	"date CHAR(8) NOT NULL DEFAULT \"\"," +
	"title VARCHAR NOT NULL DEFAULT \"\"," +
	"comment TEXT NOT NULL DEFAULT \"\"," +
	"repeat VARCHAR(128) NOT NULL DEFAULT \"\");"

var DB *sql.DB

func Init(dbFile string) error {
	var err error
	DB, err = sql.Open("sqlite", dbFile)
	if err != nil {
		return err
	}

	if _, err := DB.Exec(schema); err != nil {
		return err
	}

	return nil
}
