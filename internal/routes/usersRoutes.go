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

	userRouter := router.PathPrefix("/users").Subrouter()
	userRouter.Use(middleware.AuthMiddleware)

	// Маршрути для маніпуляцій з користувачами
	userRouter.Handle("/{uuid}", http.HandlerFunc(users.SoftDeleteUserHandler)).Methods("DELETE")
	userRouter.Handle("/{uuid}", http.HandlerFunc(users.UpdateUserHandler)).Methods("PUT")
}
