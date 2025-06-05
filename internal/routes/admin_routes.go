package routes

import (
	"github.com/gorilla/mux"

	"external-backend-go/internal/handler"
	"external-backend-go/internal/logger"
	"external-backend-go/internal/middleware"
	"external-backend-go/internal/store"
)

func setupAdminRoutes(router *mux.Router, authHandler *handler.AuthHandler, itemHandler *handler.ItemHandler, jwtSecret string, userStore store.UserStore, roleStore store.RoleStore, appLogger *logger.Logger) {
	adminRouter := router.PathPrefix("/admin").Subrouter()

	adminRouter.Use(middleware.AuthMiddleware(jwtSecret, appLogger))
	adminRouter.Use(middleware.AuthRoleMiddleware("admin", userStore, roleStore, appLogger))

	adminRouter.HandleFunc("/items", itemHandler.CreateItem).Methods("POST")
	adminRouter.HandleFunc("/items/{id}", itemHandler.UpdateItem).Methods("PUT")
	adminRouter.HandleFunc("/items/{id}", itemHandler.DeleteItem).Methods("DELETE")
	adminRouter.HandleFunc("/users/{id}/role", authHandler.UpdateUserRole).Methods("PUT")

	// adminRouter.HandleFunc("/categories", categoryHandler.CreateCategory).Methods("POST")
}
