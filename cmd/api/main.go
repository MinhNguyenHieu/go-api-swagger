package main

import (
	"log"

	"external-backend-go/configs"
	"external-backend-go/internal/app"
)

// @title Go API Application
// @version 1.0
// @description This is a sample Go API application with JWT Authentication, SQLC, and PostgreSQL.
// @termsOfService http://swagger.io/terms/

// @contact.name API Support
// @contact.url http://www.swagger.io/support
// @contact.email support@swagger.io

// @license.name Apache 2.0
// @license.url http://www.apache.org/licenses/LICENSE-2.0.html

// @host localhost:8080
// @BasePath /api/v1
// @schemes http

// @securityDefinitions.apikey ApiKeyAuth
// @in header
// @name Authorization

func main() {
	cfg := configs.LoadConfig()
	if cfg == nil {
		log.Fatal("Failed to load application configuration.")
	}

	apiApp := app.NewApp(cfg)
	apiApp.Initialize()

	apiApp.Run()
}
