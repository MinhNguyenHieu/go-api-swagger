package routes

import (
	"net/http"

	"github.com/gorilla/mux"

	"external-backend-go/internal/handler"
	"external-backend-go/internal/logger"
	"external-backend-go/internal/middleware"
	"external-backend-go/internal/utility"
)

func setupPublicRoutes(router *mux.Router, authHandler *handler.AuthHandler, basicAuthUser, basicAuthPass string, appLogger *logger.Logger) {
	router.HandleFunc("/register", authHandler.RegisterUser).Methods("POST")
	router.HandleFunc("/login", authHandler.LoginUser).Methods("POST")
	// router.HandleFunc("/verify-email", authHandler.VerifyEmail).Methods("GET")
	// router.HandleFunc("/forgot-password", authHandler.ForgotPassword).Methods("POST")
	// router.HandleFunc("/reset-password", authHandler.ResetPassword).Methods("POST")

	basicAuthRouter := router.PathPrefix("/basic-auth").Subrouter()
	basicAuthRouter.Use(middleware.BasicAuthMiddleware(basicAuthUser, basicAuthPass, func(w http.ResponseWriter, r *http.Request, err error) {
		utility.UnauthorizedBasicErrorResponse(w, r, err, appLogger)
	}))
	basicAuthRouter.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		utility.JSONResponse(w, http.StatusOK, map[string]string{"status": "Basic Auth Health OK"})
	}).Methods("GET")
}
