package users

import (
	"blogAPI/pkg/middleware"
	"encoding/json"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

func SoftDeleteUserHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	// Отримання UUID акаунту
	UUID := vars["uuid"]

	// Перевірка на існування вказаного акаунту
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

	// Перевірка чи має користувач доступ до видалення акаунту
	if UUID != middleware.User_UUID {
		http.Error(w, "You do not have the right to delete this user!", http.StatusForbidden)
		return
	}

	// Функція soft-видалення акаунту
	err = SoftDeleteUser(UUID)
	if err != nil {
		http.Error(w, "Error deleting user!", http.StatusNotFound)
		log.Println(err)
		return
	}

	// Інформування про успішне видалення акаунту
	w.WriteHeader(http.StatusNoContent)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"message": "User deleted successfully",
	})
}
