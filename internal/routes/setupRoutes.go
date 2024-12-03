package routes

import "github.com/gorilla/mux"

func SetupRoutes(router *mux.Router) {
	SetupUserRoutes(router)
	SetupCompanyRoutes(router)
	SetupArticleRoutes(router)
	SetupCategoryRoutes(router)
}
