// Package db содержит код для инициализации и работы с SQLite-базой задач.
package db

import (
	"database/sql"
	"errors"
	"fmt"
	"os"

	_ "modernc.org/sqlite"
)

const schema = `
CREATE TABLE IF NOT EXISTS scheduler (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    date CHAR(8) NOT NULL DEFAULT "",
    title VARCHAR(256) NOT NULL DEFAULT "",
    comment TEXT NOT NULL DEFAULT "",
    repeat VARCHAR(128) NOT NULL DEFAULT ""
);
CREATE INDEX IF NOT EXISTS scheduler_date_idx ON scheduler(date);
`

// DB — подключение к базе данных задач, созданное функцией Init.
var DB *sql.DB

// Init открывает SQLite-файл. Если файл ещё не существует, создаёт таблицу scheduler
// и индекс для сортировки задач по дате.
func Init(dbFile string) error {
	_, err := os.Stat(dbFile)
	install := errors.Is(err, os.ErrNotExist)
	if err != nil && !install {
		return fmt.Errorf("проверка файла базы данных: %w", err)
	}

	database, err := sql.Open("sqlite", dbFile)
	if err != nil {
		return fmt.Errorf("открытие базы данных: %w", err)
	}

	if install {
		// Файл создан SQLite при первом открытии; схему добавляем ниже.
	}
	if _, err = database.Exec(schema); err != nil {
		database.Close()
		return fmt.Errorf("создание схемы базы данных: %w", err)
	}

	DB = database
	return nil
}
