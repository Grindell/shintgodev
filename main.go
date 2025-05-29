package main

import (
	"log"
	"net/http"

	"shintgodev/db"
	"shintgodev/routs"
)

func main() {
	db.Init()

	// Один маршрут для работы со всеми пользователями
	http.HandleFunc("/users", routs.HandlerUser)

	// Отдельный маршрут для получения по ID
	http.HandleFunc("/users/", routs.HandlerUserGet)

	log.Println("Сервер запущен на :8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
