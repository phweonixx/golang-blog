package users

import (
	"blogAPI/internal/auth"
	"blogAPI/pkg/middleware"
	"encoding/json"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"golang.org/x/crypto/bcrypt"
)

func UpdateUserHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	// Отримання ID категорії
	UUID := vars["uuid"]

	// Перевірка на існування вказаного аккаунту
	var exists bool
	exists, err := checkUserExists(UUID)
	if err != nil {
		http.Error(w, "Error checking account existence:", http.StatusInternalServerError)
		return
	}
	if !exists {
		http.Error(w, "Account with this UUID not found!", http.StatusNotFound)
		return
	}

	// Перевірка чи має користувач доступ до оновлення акаунту
	if UUID != middleware.User_UUID {
		http.Error(w, "You do not have the right to change this user!", http.StatusForbidden)
		return
	}

	// Отримання тіла запросу в JSON
	var user auth.User
	err = json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		http.Error(w, "Invalid Input!", http.StatusBadRequest)
		log.Println("Error decoding request body:", err)
		return
	}

	exists, err = CheckUserUsernameOrEmailExists(user.Username, user.Email)
	if err != nil {
		http.Error(w, "Error checking username or email existanse", http.StatusInternalServerError)
		log.Println(err)
		return
	}
	if exists {
		http.Error(w, "user with this username or email already exists", http.StatusInternalServerError)
		return
	}

	if user.Password != "" {
		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
		if err != nil {
			return
		}
		user.Password = string(hashedPassword)
	}

	// Функція оновлення користувача
	err = UpdateUser(user, UUID)
	if err != nil {
		http.Error(w, "Error updating user!", http.StatusNotFound)
		log.Println(err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{
		"message": "User updated successfully",
	})
}
