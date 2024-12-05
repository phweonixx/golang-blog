package routes

import (
	"blogAPI/internal/handlers/companies"
	"blogAPI/pkg/middleware"
	"net/http"

	"github.com/gorilla/mux"
)

func SetupCompanyRoutes(router *mux.Router) {
	companyRouter := router.PathPrefix("/companies").Subrouter()
	companyRouter.Use(middleware.AuthMiddleware)

	// Маршрути для маніпуляцій з компаніями
	companyRouter.Handle("", http.HandlerFunc(companies.CreateCompanyHandler)).Methods("POST")
	companyRouter.Handle("", http.HandlerFunc(companies.ReadCompaniesHandler)).Methods("GET")
	companyRouter.Handle("/{uuid}", http.HandlerFunc(companies.UpdateCompanyHandler)).Methods("PUT")
	companyRouter.Handle("/{uuid}", http.HandlerFunc(companies.SoftDeleteCompanyHandler)).Methods("DELETE")
}
