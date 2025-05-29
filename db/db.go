package db

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"strings"

	_ "github.com/lib/pq"
)

var DB *sql.DB

type User struct {
	ID   int
	Name string
	Age  int
}

func Init() {
	host := os.Getenv("DB_HOST")
	port := os.Getenv("DB_PORT")
	user := os.Getenv("DB_USER")
	password := os.Getenv("DB_PASSWORD")
	dbname := os.Getenv("DB_NAME")

	connStr := fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		host, port, user, password, dbname,
	)

	var err error
	DB, err = sql.Open("postgres", connStr)
	if err != nil {
		log.Fatalf("Ошибка подключения к БД: %v", err)
	}

	// Создаем таблицу, если не существует
	createTableUsers := `
	CREATE TABLE IF NOT EXISTS users (
		id SERIAL PRIMARY KEY,
		name TEXT NOT NULL,
		age INTEGER
	);`

	if _, err := DB.Exec(createTableUsers); err != nil {
		log.Fatalf("Ошибка при создании таблицы: %v", err)
	} else {
		fmt.Println("Таблица users создана или уже существует")
	}
}

func CreateUser(user User) (User, error) {
	query := "INSERT INTO users (name, age) VALUES ($1, $2) RETURNING id"
	err := DB.QueryRow(query, user.Name, user.Age).Scan(&user.ID)
	if err != nil {
		log.Printf("Ошибка при добавлении пользователя: %v", err)
		return User{}, err
	}
	return user, nil
}

func GetUserByID(id int) (User, error) {
	query := "SELECT id, name, age FROM users WHERE id = $1"

	var u User
	err := DB.QueryRow(query, id).Scan(&u.ID, &u.Name, &u.Age)
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
		if err := rows.Scan(&u.ID, &u.Name, &u.Age); err != nil {
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
	query := "UPDATE users SET name = $1, age = $2 WHERE id = $3"
	_, err := DB.Exec(query, user.Name, user.Age, user.ID)
	if err != nil {
		return User{}, err
	}
	return user, nil
}

func PatchUser(user User) (User, error) {
	setParts := []string{}
	args := []interface{}{}
	argPos := 1

	if user.Name != "" {
		setParts = append(setParts, fmt.Sprintf("name = $%d", argPos))
		args = append(args, user.Name)
		argPos++
	}
	if user.Age != 0 {
		setParts = append(setParts, fmt.Sprintf("age = $%d", argPos))
		args = append(args, user.Age)
		argPos++
	}

	if len(setParts) == 0 {
		return user, nil
	}

	query := fmt.Sprintf("UPDATE users SET %s WHERE id = $%d", strings.Join(setParts, ", "), argPos)
	args = append(args, user.ID)

	_, err := DB.Exec(query, args...)
	if err != nil {
		return User{}, err
	}

	return user, nil
}

func DeleteUser(userID int) error {
	query := "DELETE FROM users WHERE id = $1"
	_, err := DB.Exec(query, userID)
	return err
}
