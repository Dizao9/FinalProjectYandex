package database

import (
	"database/sql"
	"fmt"
	"os"

	_ "modernc.org/sqlite"
)

var DB *sql.DB

const Schema = `CREATE TABLE IF NOT EXISTS scheduler (
	id INTEGER PRIMARY KEY AUTOINCREMENT,
	date CHAR(8) NOT NULL DEFAULT "",
	comment TEXT,
	title VARCHAR(256),
	repeat VARCHAR(128)
);

CREATE INDEX scheduler_date ON scheduler(date);
`

func Init(dbFile string) error {
	_, err := os.Stat(dbFile)
	var install bool
	if err != nil {
		install = true
	}

	DB, err = sql.Open("sqlite", dbFile)

	if err != nil {
		return fmt.Errorf("Ошибка при открытии базы данных: %v", err)
	}

	if install {
		_, err := DB.Exec(Schema)
		if err != nil {
			return fmt.Errorf("Ошибка при создании таблицы schedular: %v", err)
		}
	}
	return nil
}
