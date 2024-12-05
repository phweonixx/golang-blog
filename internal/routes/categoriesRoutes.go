package routes

import (
	"blogAPI/internal/handlers/categories"
	"blogAPI/pkg/middleware"
	"net/http"

	"github.com/gorilla/mux"
)

func SetupCategoryRoutes(router *mux.Router) {
	categoryRouter := router.PathPrefix("/categories").Subrouter()
	categoryRouter.Use(middleware.AuthMiddleware)

	// Маршрути для маніпуляцій з категоріями
	categoryRouter.Handle("", http.HandlerFunc(categories.CreateCategory)).Methods("POST")
	categoryRouter.Handle("", http.HandlerFunc(categories.ReadCategories)).Methods("GET")

	categoryRouter.Handle("/search", http.HandlerFunc(categories.SearchCategories)).Methods("GET")

	categoryRouter.Handle("/my-categories", http.HandlerFunc(categories.ReadCategoriesMy)).Methods("GET")

	categoryRouter.Handle("/{id}", http.HandlerFunc(categories.ReadCategory)).Methods("GET")
	categoryRouter.Handle("/{id}", http.HandlerFunc(categories.UpdateCategory)).Methods("PUT")
	categoryRouter.Handle("/{id}", http.HandlerFunc(categories.DeleteCategory)).Methods("DELETE")
}
