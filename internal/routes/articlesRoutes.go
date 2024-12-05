package routes

import (
	"blogAPI/internal/handlers/articles"
	"blogAPI/pkg/middleware"
	"net/http"

	"github.com/gorilla/mux"
)

func SetupArticleRoutes(router *mux.Router) {
	articleRouter := router.PathPrefix("/articles").Subrouter()
	articleRouter.Use(middleware.AuthMiddleware)

	// Маршрути для маніпуляцій зі статтями
	articleRouter.Handle("", http.HandlerFunc(articles.CreateArticle)).Methods("POST")
	articleRouter.Handle("", http.HandlerFunc(articles.ReadArticles)).Methods("GET")

	articleRouter.Handle("/search", http.HandlerFunc(articles.SearchArticles)).Methods("GET")

	articleRouter.Handle("/my-articles", http.HandlerFunc(articles.ReadArticlesMy)).Methods("GET")

	articleRouter.Handle("/{id}", http.HandlerFunc(articles.ReadArticle)).Methods("GET")
	articleRouter.Handle("/{id}", http.HandlerFunc(articles.UpdateArticle)).Methods("PUT")
	articleRouter.Handle("/{id}", http.HandlerFunc(articles.DeleteArticle)).Methods("DELETE")
}
