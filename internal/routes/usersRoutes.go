package routes

import (
	"blogAPI/internal/auth"
	"blogAPI/internal/handlers/users"
	"blogAPI/pkg/middleware"
	"net/http"

	"github.com/gorilla/mux"
)

func SetupUserRoutes(router *mux.Router) {
	// Маршрут реєстрації користувача
	router.HandleFunc("/register", auth.RegisterHandler).Methods("POST")

	// Маршрут для входу в акаунт
	router.HandleFunc("/login", auth.LoginHandler).Methods("POST")

	// Маршрути для маніпуляцій з користувачами
	router.Handle("/users/{uuid}", middleware.AuthMiddleware(http.HandlerFunc(users.SoftDeleteUserHandler))).Methods("DELETE")
	router.Handle("/users/{uuid}", middleware.AuthMiddleware(http.HandlerFunc(users.UpdateUserHandler))).Methods("PUT")
}
