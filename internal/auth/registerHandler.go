package auth

import (
	"blogAPI/internal/helpers"
	"blogAPI/internal/models"
	"encoding/json"
	"log"
	"net/http"
)

func RegisterHandler(w http.ResponseWriter, r *http.Request) {
	// Отримання тіла запросу в JSON
	var user models.User
	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		http.Error(w, "Invalid input!", http.StatusBadRequest)
		log.Println(err)
		return
	}

	// Функція реєстрації користувача та видачі токену
	err = RegisterUser(&user)
	if err != nil {
		switch err {
		case ErrUsernameTooShort:
			http.Error(w, "The username length must be more than 3 characters!", http.StatusBadRequest)
		case ErrPasswordTooShort:
			http.Error(w, "The password length must be more than 8 characters!", http.StatusBadRequest)
		case ErrFirstNameTooShort:
			http.Error(w, "First name length must be more than 2 characters!", http.StatusBadRequest)
		case ErrLastNameTooShort:
			http.Error(w, "Last name length must be more than 2 characters!", http.StatusBadRequest)
		case ErrEmailTooShort:
			http.Error(w, "Email length must be more than 6 characters!", http.StatusBadRequest)
		case ErrEmailOrUsernameUsed:
			http.Error(w, "This Email or Username is already used!", http.StatusConflict)
		case ErrEmailInvalidFormat:
			http.Error(w, "Invalid Email format!", http.StatusBadRequest)
		default:
			http.Error(w, "Error during registration!", http.StatusInternalServerError)
			log.Println(err)
		}
		return
	}

	// Інформування про успішне створення акаунту
	helpers.SendJSONResponse(w, http.StatusCreated, "Account created successfully", user)
}
