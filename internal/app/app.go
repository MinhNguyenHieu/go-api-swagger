package app

import (
	"database/sql"
	"log"
	"net/http"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/gorilla/mux"

	"external-backend-go/configs"
	"external-backend-go/db/sqlc"
	"external-backend-go/internal/database"
	"external-backend-go/internal/elasticsearch"
	"external-backend-go/internal/email"
	"external-backend-go/internal/handler"
	"external-backend-go/internal/logger"
	"external-backend-go/internal/middleware"
	"external-backend-go/internal/routes"
	"external-backend-go/internal/service"
	"external-backend-go/internal/store"
)

type App struct {
	Config                  *configs.Config
	Router                  *mux.Router
	DB                      *sql.DB
	Queries                 *sqlc.Queries
	UserStore               store.UserStore
	ItemStore               store.ItemStore
	RoleStore               store.RoleStore
	PasswordResetTokenStore store.PasswordResetTokenStore
	SessionStore            store.SessionStore
	SearchStore             store.SearchStore
	ElasticsearchClient     *elasticsearch.Client

	AuthService *service.AuthService
	ItemService *service.ItemService
	AuthHandler *handler.AuthHandler
	ItemHandler *handler.ItemHandler
	EmailSender email.EmailSender
	RateLimiter *middleware.RateLimiter
	Logger      *logger.Logger
	Validator   *validator.Validate
}

func NewApp(cfg *configs.Config) *App {
	return &App{
		Config: cfg,
		Router: mux.NewRouter(),
		Logger: logger.NewLogger(),
	}
}

func (a *App) Initialize() {
	var err error
	maxRetries := 10
	for i := 0; i < maxRetries; i++ {
		a.DB, err = database.ConnectDB(a.Config.DatabaseURL)
		if err == nil {
			err = a.DB.Ping()
			if err == nil {
				a.Logger.Info("Successfully connected to database!")
				break
			}
		}
		a.Logger.Error("Failed to connect or initialize database: %v. Retrying in 5 seconds...", err)
		time.Sleep(5 * time.Second)
		if i == maxRetries-1 {
			a.Logger.Fatal("Exceeded max retries for database connection: %v", err)
		}
	}

	err = database.RunMigrations(a.Config.DatabaseURL, "db/migrations")
	if err != nil {
		a.Logger.Fatal("Failed to run migrations: %v", err)
	}
	a.Logger.Info("Database migrations applied successfully.")

	a.Queries = sqlc.New(a.DB)

	baseRepo := store.NewBaseRepository(a.DB, a.Logger)

	// Initialize all stores
	a.UserStore = store.NewUserStore(a.DB, a.Queries, baseRepo)
	a.ItemStore = store.NewItemStore(a.DB, a.Queries, baseRepo)
	a.RoleStore = store.NewRoleStore(a.DB, a.Queries, baseRepo)
	a.PasswordResetTokenStore = store.NewPasswordResetTokenStore(a.DB, a.Queries, baseRepo)
	a.SessionStore = store.NewSessionStore(a.DB, a.Queries, baseRepo)

	a.ElasticsearchClient, err = elasticsearch.NewElasticsearchClient("http://elasticsearch:9200")
	if err != nil {
		a.Logger.Fatal("Failed to connect to Elasticsearch: %v", err)
	}
	a.Logger.Info("Successfully connected to Elasticsearch!")

	a.SearchStore = store.NewSearchStore(a.ElasticsearchClient, a.Logger)

	a.EmailSender = email.NewSMTPEmailSender(
		a.Config.SMTP.Host,
		a.Config.SMTP.Port,
		a.Config.SMTP.User,
		a.Config.SMTP.Pass,
		a.Config.SMTP.SenderEmail,
	)

	a.Validator = validator.New()

	// Initialize services
	a.AuthService = service.NewAuthService(a.UserStore, a.RoleStore, a.SessionStore, a.PasswordResetTokenStore, a.Config.JWTSecret, a.EmailSender)
	a.ItemService = service.NewItemService(a.ItemStore, a.SearchStore)

	// Initialize handlers, passing logger and validator
	a.ItemHandler = handler.NewItemHandler(a.ItemService, a.Logger, a.Validator)
	a.AuthHandler = handler.NewAuthHandler(a.AuthService, a.Logger, a.Validator)

	// Initialize Rate Limiter
	a.RateLimiter = middleware.NewRateLimiter(
		a.Config.RateLimiter.Enabled,
		a.Config.RateLimiter.RPS,
		a.Config.RateLimiter.Burst,
		a.Config.RateLimiter.TTL,
	)
	a.Logger.Info("Rate Limiter initialized. Enabled: %t, RPS: %.2f, Burst: %d", a.Config.RateLimiter.Enabled, a.Config.RateLimiter.RPS, a.Config.RateLimiter.Burst)

	routes.SetupAPIRoutes(routes.AppDependencies{
		Router:        a.Router,
		AuthHandler:   a.AuthHandler,
		ItemHandler:   a.ItemHandler,
		JWTSecret:     a.Config.JWTSecret,
		UserStore:     a.UserStore,
		RoleStore:     a.RoleStore,
		RateLimiter:   a.RateLimiter,
		BasicAuthUser: a.Config.Auth.Basic.User,
		BasicAuthPass: a.Config.Auth.Basic.Pass,
		AppLogger:     a.Logger,
		SearchStore:   a.SearchStore,
		// Validator:     a.Validator,
	})

	a.Logger.Info("Application initialization complete.")
}

func (a *App) Run() {
	addr := ":" + a.Config.AppPort
	a.Logger.Info("Server is running on %s", addr)
	log.Fatal(http.ListenAndServe(addr, a.Router))
}
