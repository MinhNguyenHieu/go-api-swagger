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

type AppDependencies struct {
	Router        *mux.Router
	AuthHandler   *handler.AuthHandler
	ItemHandler   *handler.ItemHandler
	JWTSecret     string
	UserStore     store.UserStore
	RoleStore     store.RoleStore
	RateLimiter   *middleware.RateLimiter
	BasicAuthUser string
	BasicAuthPass string
	AppLogger     *logger.Logger
	SearchStore   store.SearchStore
}

func SetupAPIRoutes(deps AppDependencies) {
	corsMiddleware := middleware.NewCORS()
	deps.Router.Use(corsMiddleware.Handler)

	deps.Router.Use(middleware.LoggerMiddleware)

	apiV1Router := deps.Router.PathPrefix("/api/v1").Subrouter()

	apiV1Router.Use(middleware.RateLimiterMiddleware(deps.RateLimiter, func(w http.ResponseWriter, r *http.Request, retryAfter string) {
		utility.RateLimitExceededResponse(w, r, retryAfter, deps.AppLogger)
	}))

	setupPublicRoutes(
		apiV1Router,
		deps.AuthHandler,
		deps.BasicAuthUser,
		deps.BasicAuthPass,
		deps.AppLogger,
	)

	setupProtectedRoutes(
		apiV1Router,
		deps.AuthHandler,
		deps.ItemHandler,
		deps.JWTSecret,
		deps.AppLogger,
	)

	setupAdminRoutes(
		apiV1Router,
		deps.AuthHandler,
		deps.ItemHandler,
		deps.JWTSecret,
		deps.UserStore,
		deps.RoleStore,
		deps.AppLogger,
	)

	deps.Router.PathPrefix("/swagger/").Handler(httpSwagger.WrapHandler)
}
