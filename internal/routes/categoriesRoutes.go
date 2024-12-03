package routes

import (
	"blogAPI/internal/handlers/categories"
	"blogAPI/pkg/middleware"
	"net/http"

	"github.com/gorilla/mux"
)

func SetupCategoryRoutes(router *mux.Router) {
	// Маршрути для маніпуляцій з категоріями
	router.Handle("/categories", middleware.AuthMiddleware(http.HandlerFunc(categories.CreateCategory))).Methods("POST")
	router.Handle("/categories", middleware.AuthMiddleware(http.HandlerFunc(categories.ReadCategories))).Methods("GET")

	router.Handle("/categories/search", middleware.AuthMiddleware(http.HandlerFunc(categories.SearchCategories))).Methods("GET")

	router.Handle("/categories/my-categories", middleware.AuthMiddleware(http.HandlerFunc(categories.ReadCategoriesMy))).Methods("GET")

	router.Handle("/categories/{id}", middleware.AuthMiddleware(http.HandlerFunc(categories.ReadCategory))).Methods("GET")
	router.Handle("/categories/{id}", middleware.AuthMiddleware(http.HandlerFunc(categories.UpdateCategory))).Methods("PUT")
	router.Handle("/categories/{id}", middleware.AuthMiddleware(http.HandlerFunc(categories.DeleteCategory))).Methods("DELETE")
}
