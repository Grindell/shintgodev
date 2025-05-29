package db

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/mattn/go-sqlite3"
)

var DB *sql.DB

func Init() {
	var err error
	DB, err = sql.Open("sqlite3", "data.db")
	if err != nil {
		log.Fatalf("Ошибка подключения к БД: %v", err)
	}

	// Пример создания таблицы
	createTableUsers := `
		CREATE TABLE IF NOT EXISTS users (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			name TEXT NOT NULL,
			age INTEGER
		);`

	if _, err := DB.Exec(createTableUsers); err != nil {
		log.Fatalf("Ошибка при создании таблицы: %v", err)
	} else {
		fmt.Println("База данных users создана:")
	}

}

func CreateUser(user users.User) error {
	query := "INSERT INTO users (name, age) VALUES (?, ?)"
	_, err := DB.Exec(query, user.Name, user.Age)
	if err != nil {
		log.Printf("Ошибка при добавлении пользователя: %v", err)
		return err
	}
	return nil
}
