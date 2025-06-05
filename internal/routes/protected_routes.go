package routes

import (
	"github.com/gorilla/mux"

	"external-backend-go/internal/handler"
	"external-backend-go/internal/logger"
	"external-backend-go/internal/middleware"
)

func setupProtectedRoutes(router *mux.Router, authHandler *handler.AuthHandler, itemHandler *handler.ItemHandler, jwtSecret string, appLogger *logger.Logger) {
	protectedRouter := router.PathPrefix("").Subrouter()
	protectedRouter.Use(middleware.AuthMiddleware(jwtSecret, appLogger))

	protectedRouter.HandleFunc("/protected", authHandler.ProtectedEndpoint).Methods("GET")

	protectedRouter.HandleFunc("/items", itemHandler.GetItems).Methods("GET")
	protectedRouter.HandleFunc("/items/{id}", itemHandler.GetItem).Methods("GET")

	// protectedRouter.HandleFunc("/profile", userHandler.GetUserProfile).Methods("GET")
}
