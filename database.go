package main

import (
	"database/sql"
	"log"

	_ "github.com/mattn/go-sqlite3"
)

var db *sql.DB

// Инициализация базы данных
func initDB() {
	var err error
	db, err = sql.Open("sqlite3", "./logs.db")
	if err != nil {
		log.Fatal("Ошибка подключения к базе данных:", err)
	}

	_, err = db.Exec(`
        CREATE TABLE IF NOT EXISTS logs (
            id INTEGER PRIMARY KEY AUTOINCREMENT,
            tab TEXT,
            status TEXT,
            data TEXT,
            timestamp DATETIME DEFAULT CURRENT_TIMESTAMP
        )`)
	if err != nil {
		log.Fatal("Ошибка при создании таблицы:", err)
	}
}

// Функция для добавления записи в базу данных
func insertLog(tab, status, data string) error {
	_, err := db.Exec("INSERT INTO logs (tab, status, data) VALUES (?, ?, ?)", tab, status, data)
	return err
}
