package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/hadi-projects/xyz-finance-go/config"
	"github.com/hadi-projects/xyz-finance-go/internal/entity"
	"github.com/hadi-projects/xyz-finance-go/internal/handler"
	"github.com/hadi-projects/xyz-finance-go/internal/repository"
	"github.com/hadi-projects/xyz-finance-go/internal/router"
	services "github.com/hadi-projects/xyz-finance-go/internal/service"
	"github.com/hadi-projects/xyz-finance-go/pkg/database"
	"gorm.io/gorm"
)

type Application struct {
	Config *config.Config
	DB     *gorm.DB
	Router *gin.Engine
	Server *http.Server
}

func main() {
	app := &Application{}
	app.initializeConfig()
	app.initializeDatabase()
	app.setupRouter()
	app.run()
}

// initializeConfig loads the application configuration
func (app *Application) initializeConfig() {
	cfg, err := config.NewConfig()
	if err != nil {
		log.Fatalf("Failed to load configuration: %v", err)
	}
	app.Config = cfg
	log.Println("Configuration loaded successfully")
}

// initializeDatabase connects to the database and runs migrations
func (app *Application) initializeDatabase() {
	db, err := database.NewMySQLConnection(app.Config)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	app.DB = db

	// Auto-migrate database tables
	if err := db.AutoMigrate(
		&entity.Role{},
		&entity.Permission{},
		&entity.User{},
		&entity.RefreshToken{},
		&entity.TenorLimit{},
	); err != nil {
		log.Fatalf("Failed to migrate database: %v", err)
	}
	log.Println("Database migration completed successfully")

	database.SeedRBAC(app.DB)
	database.SeedUser(app.DB)

}

// setupRouter initializes all dependencies and configures routes
func (app *Application) setupRouter() {

	userRepo := repository.NewUserRepository(app.DB)
	authService := services.NewAuthService(userRepo)
	refreshTokenRepo := repository.NewRefreshTokenRepository(app.DB)
	jwtService := services.NewJWTService(app.Config.AppPort, app.Config.JWT.ExpiryHours, refreshTokenRepo)
	authHandler := handler.NewAuthHandler(authService, jwtService)

	limitRepo := repository.NewLimitRepository(app.DB)
	limitService := services.NewLimitService(limitRepo)
	limitHandler := handler.NewLimitHandler(limitService)

	appRouter := router.NewRouter(app.Config, authHandler, limitHandler)
	app.Router = appRouter.SetupRoutes()

	log.Println("Router configured successfully")
}

// run starts the HTTP server and handles graceful shutdown
func (app *Application) run() {
	app.Server = &http.Server{
		Addr:           ":" + app.Config.AppPort,
		Handler:        app.Router,
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   10 * time.Second,
		MaxHeaderBytes: 1 << 20, // 1 MB
	}

	// Start server in goroutine
	go func() {
		fmt.Printf("🚀 Server running on port %s\n", app.Config.AppPort)
		if err := app.Server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("❌ Server failed to start: %v", err)
		}
	}()

	// Wait for interrupt signal for graceful shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Println("Shutting down server...")

	// Graceful shutdown with 5 second timeout
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := app.Server.Shutdown(ctx); err != nil {
		log.Fatalf("Server forced to shutdown: %v", err)
	}

	log.Println("Server exited successfully")
}
