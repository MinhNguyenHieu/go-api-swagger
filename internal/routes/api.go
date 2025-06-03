package routes

import (
	"net/http"

	"github.com/gorilla/mux"
	httpSwagger "github.com/swaggo/http-swagger"

	_ "external-backend-go/docs"
	"external-backend-go/internal/handler"
	"external-backend-go/internal/logger"
	"external-backend-go/internal/middleware"
	"external-backend-go/internal/store"
	"external-backend-go/internal/utility"
)

// SetupAPIRoutes configures all API endpoints for the application.
func SetupAPIRoutes(
	router *mux.Router,
	authHandler *handler.AuthHandler,
	itemHandler *handler.ItemHandler,
	jwtSecret string,
	userStore store.UserStore,
	rateLimiter *middleware.RateLimiter,
	basicAuthUser string,
	basicAuthPass string,
	appLogger *logger.Logger,
	roleStore store.RoleStore,
) {
	// Apply LoggerMiddleware to all requests
	router.Use(middleware.LoggerMiddleware)

	apiV1Router := router.PathPrefix("/api/v1").Subrouter()

	// Apply RateLimiterMiddleware to all API v1 routes
	apiV1Router.Use(middleware.RateLimiterMiddleware(rateLimiter, func(w http.ResponseWriter, r *http.Request, retryAfter string) {
		utility.RateLimitExceededResponse(w, r, retryAfter, appLogger)
	}))

	// Public routes (no authentication needed)
	apiV1Router.HandleFunc("/register", authHandler.RegisterUser).Methods("POST")
	apiV1Router.HandleFunc("/login", authHandler.LoginUser).Methods("POST")

	// Protected routes (require JWT authentication)
	protectedRouter := apiV1Router.PathPrefix("/").Subrouter()
	protectedRouter.Use(middleware.AuthMiddleware(jwtSecret, appLogger))

	protectedRouter.HandleFunc("/protected", authHandler.ProtectedEndpoint).Methods("GET")

	// Item CRUD operations (protected)
	protectedRouter.HandleFunc("/items", itemHandler.CreateItem).Methods("POST")
	protectedRouter.HandleFunc("/items/{id}", itemHandler.GetItem).Methods("GET")
	protectedRouter.HandleFunc("/items/{id}", itemHandler.UpdateItem).Methods("PUT")
	protectedRouter.HandleFunc("/items/{id}", itemHandler.DeleteItem).Methods("DELETE")
	protectedRouter.HandleFunc("/items", itemHandler.GetItems).Methods("GET")

	// Admin routes (require JWT authentication and 'admin' role)
	adminRouter := apiV1Router.PathPrefix("/admin").Subrouter()
	adminRouter.Use(middleware.AuthMiddleware(jwtSecret, appLogger))
	adminRouter.Use(middleware.AuthRoleMiddleware("admin", userStore, roleStore, appLogger))

	adminRouter.HandleFunc("/users/{id}/role", authHandler.UpdateUserRole).Methods("PUT")

	// Basic Auth example route
	basicAuthRouter := apiV1Router.PathPrefix("/basic-auth").Subrouter()
	basicAuthRouter.Use(middleware.BasicAuthMiddleware(basicAuthUser, basicAuthPass, func(w http.ResponseWriter, r *http.Request, err error) {
		utility.UnauthorizedBasicErrorResponse(w, r, err, appLogger)
	}))
	basicAuthRouter.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		utility.JSONResponse(w, http.StatusOK, map[string]string{"status": "Basic Auth Health OK"})
	}).Methods("GET")

	// Swagger UI endpoint
	router.PathPrefix("/swagger/").Handler(httpSwagger.WrapHandler)
}
