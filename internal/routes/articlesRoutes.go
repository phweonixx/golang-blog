package routes

import (
	"blogAPI/internal/handlers/articles"
	"blogAPI/pkg/middleware"
	"net/http"

	"github.com/gorilla/mux"
)

func SetupArticleRoutes(router *mux.Router) {
	// Маршрути для маніпуляцій зі статтями
	router.Handle("/articles", middleware.AuthMiddleware(http.HandlerFunc(articles.CreateArticle))).Methods("POST")
	router.Handle("/articles", middleware.AuthMiddleware(http.HandlerFunc(articles.ReadArticles))).Methods("GET")

	router.Handle("/articles/search", middleware.AuthMiddleware(http.HandlerFunc(articles.SearchArticles))).Methods("GET")

	router.Handle("/articles/my-articles", middleware.AuthMiddleware(http.HandlerFunc(articles.ReadArticlesMy))).Methods("GET")

	router.Handle("/articles/{id}", middleware.AuthMiddleware(http.HandlerFunc(articles.ReadArticle))).Methods("GET")
	router.Handle("/articles/{id}", middleware.AuthMiddleware(http.HandlerFunc(articles.UpdateArticle))).Methods("PUT")
	router.Handle("/articles/{id}", middleware.AuthMiddleware(http.HandlerFunc(articles.DeleteArticle))).Methods("DELETE")
}
