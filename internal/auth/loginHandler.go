package auth

import (
	"blogAPI/internal/helpers"
	"blogAPI/internal/models"
	"encoding/json"
	"log"
	"net/http"
)

func LoginHandler(w http.ResponseWriter, r *http.Request) {
	// Отримання тіла запросу в JSON
	var credentials models.Credentials
	err := json.NewDecoder(r.Body).Decode(&credentials)
	if err != nil {
		http.Error(w, "Invalid Input!", http.StatusBadRequest)
		log.Println("Error decoding request body:", err)
		return
	}

	// Перевірка даних на валідність
	if credentials.Username == "" && credentials.Email == "" {
		http.Error(w, "Username or Email is required!", http.StatusBadRequest)
		return
	}
	if credentials.Password == "" {
		http.Error(w, "Password is required!", http.StatusBadRequest)
		return
	}

	// Перевірка на існування вказаного акаунту
	exists, err := helpers.CheckUserExistsLogin(credentials.Username, credentials.Email)
	if err != nil {
		http.Error(w, "Error checking user", http.StatusInternalServerError)
		log.Println("Error checking user existence:", err)
		return
	}
	if !exists {
		http.Error(w, "Invalid credentials", http.StatusUnauthorized)
		return
	}

	// Функція створення токену
	token, err := LoginUser(&credentials)
	if err != nil {
		http.Error(w, "Login error", http.StatusBadRequest)
		log.Println(err)
		return
	}

	// Інформування про успішний вхід в акаунт
	helpers.SendJSONResponse(w, http.StatusOK, "Login successfull", token)
}
