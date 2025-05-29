package routs

import (
	"encoding/json"
	"net/http"
	"shintgodev/db"
	"strconv"
	"strings"
)

var handlers = map[string]http.HandlerFunc{
	"POST": HandlerUserPost,
	"GET":  HandlerUserGet,
	// "PUT":    HandlerUserPut,
	// "PATCH":  HandlerUserPatch,
	// "DELETE": HandlerUserDelete,
}

func HandlerUser(w http.ResponseWriter, r *http.Request) {
	if handler, ok := handlers[r.Method]; ok {
		handler(w, r)
	} else {
		http.Error(w, "Метод не поддерживается", http.StatusMethodNotAllowed)
	}
}

func HandlerUserPost(w http.ResponseWriter, r *http.Request) {
	var user db.User
	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		http.Error(w, "Ошибка разбора JSON: "+err.Error(), http.StatusBadRequest)
		return
	}

	createdUser, err := db.CreateUser(user)
	if err != nil {
		http.Error(w, "Ошибка при создании пользователя: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(createdUser)
}

func HandlerUserGet(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Метод не поддерживается", http.StatusMethodNotAllowed)
		return
	}

	parts := strings.Split(strings.Trim(r.URL.Path, "/"), "/")

	// Если путь просто /users — вернуть всех
	if len(parts) == 1 {
		users, err := db.GetAllUsers()
		if err != nil {
			http.Error(w, "Ошибка при получении пользователей", http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(users)
		return
	}

	// Если путь /users/{id} — вернуть конкретного пользователя
	if len(parts) == 2 {
		id, err := strconv.Atoi(parts[1])
		if err != nil {
			http.Error(w, "Некорректный ID", http.StatusBadRequest)
			return
		}

		user, err := db.GetUserByID(id)
		if err != nil {
			http.Error(w, "Пользователь не найден", http.StatusNotFound)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(user)
		return
	}

	// Всё остальное — ошибка
	http.Error(w, "Некорректный путь", http.StatusBadRequest)
}

// func HandlerUserPut(w http.ResponseWriter, r *http.Request)

// func HandlerUserPatch(w http.ResponseWriter, r *http.Request)

// func HandlerUserAll(w http.ResponseWriter, r *http.Request)

// func HandlerUserDelete(w http.ResponseWriter, r *http.Request) {

// }
