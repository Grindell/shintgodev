package db

import (
	"database/sql"
	"fmt"
	"log"
	"strings"

	_ "github.com/mattn/go-sqlite3"
)

var DB *sql.DB

type User struct {
	ID   int
	Name string
	Age  int
}

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

func CreateUser(user User) (User, error) {
	query := "INSERT INTO users (name, age) VALUES (?, ?)"
	result, err := DB.Exec(query, user.Name, user.Age)
	if err != nil {
		log.Printf("Ошибка при добавлении пользователя: %v", err)
		return User{}, err
	}

	id, err := result.LastInsertId()
	if err != nil {
		log.Printf("Ошибка при получении ID нового пользователя: %v", err)
		return User{}, err
	}

	user.ID = int(id) // присваиваем ID из базы в структуру

	return user, nil
}

func GetUserByID(id int) (User, error) {
	query := "SELECT id, name, age FROM users WHERE id = ?"

	row := DB.QueryRow(query, id)

	var u User
	err := row.Scan(&u.ID, &u.Name, &u.Age)
	if err != nil {
		return User{}, err
	}

	return u, nil
}

func GetAllUsers() ([]User, error) {
	query := "SELECT id, name, age FROM users"

	rows, err := DB.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var users []User

	for rows.Next() {
		var u User
		err := rows.Scan(&u.ID, &u.Name, &u.Age)
		if err != nil {
			return nil, err
		}
		users = append(users, u)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return users, nil
}

func UpdateUser(user User) (User, error) {
	query := "UPDATE users SET name = ?, age = ? WHERE id = ?"

	_, err := DB.Exec(query, user.Name, user.Age, user.ID)
	if err != nil {
		return User{}, err
	}

	return user, nil
}

func PatchUser(user User) (User, error) {
	// Слайсы для частей запроса и аргументов
	setParts := []string{}
	args := []interface{}{}

	if user.Name != "" {
		setParts = append(setParts, "name = ?")
		args = append(args, user.Name)
	}
	if user.Age != 0 {
		setParts = append(setParts, "age = ?")
		args = append(args, user.Age)
	}

	if len(setParts) == 0 {
		// Нет полей для обновления
		return user, nil
	}

	// Формируем запрос
	query := "UPDATE users SET " + strings.Join(setParts, ", ") + " WHERE id = ?"
	args = append(args, user.ID)

	_, err := DB.Exec(query, args...)
	if err != nil {
		return User{}, err
	}

	return user, nil
}

func DeleteUser(userID int) error {
	query := "DELETE FROM users WHERE id = ?"
	_, err := DB.Exec(query, userID)
	if err != nil {
		return err
	}
	return nil
}
