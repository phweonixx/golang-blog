package routes

import (
	"blogAPI/internal/handlers/companies"
	"blogAPI/pkg/middleware"
	"net/http"

	"github.com/gorilla/mux"
)

func SetupCompanyRoutes(router *mux.Router) {
	// Маршрути для маніпуляцій з компаніями
	router.Handle("/companies", middleware.AuthMiddleware(http.HandlerFunc(companies.CreateCompanyHandler))).Methods("POST")
	router.Handle("/companies", middleware.AuthMiddleware(http.HandlerFunc(companies.ReadCompaniesHandler))).Methods("GET")
	router.Handle("/companies/{uuid}", middleware.AuthMiddleware(http.HandlerFunc(companies.UpdateCompanyHandler))).Methods("PUT")
	router.Handle("/companies/{uuid}", middleware.AuthMiddleware(http.HandlerFunc(companies.SoftDeleteCompanyHandler))).Methods("DELETE")
}
