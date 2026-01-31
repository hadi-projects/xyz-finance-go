package main

import (
	"context"
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
	"github.com/hadi-projects/xyz-finance-go/pkg/logger"
	"gorm.io/gorm"
)

type Application struct {
	Config *config.AppConfig
	DB     *gorm.DB
	Router *gin.Engine
	Server *http.Server
}

func main() {
	app := &Application{}
	app.initializeConfig()

	// Initialize Logger
	logger.Init(logger.Config{
		LogDir:      "logs",
		Environment: app.Config.AppEnv,
	})
	logger.SystemLogger.Info().Msg("Logger initialized")

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
	// Using std log here because logger is not init yet
	log.Println("Configuration loaded successfully")
}

// initializeDatabase connects to the database and runs migrations
func (app *Application) initializeDatabase() {
	db, err := database.NewMySQLConnection(app.Config)
	if err != nil {
		logger.SystemLogger.Fatal().Err(err).Msg("Failed to connect to database")
	}
	app.DB = db

	// Auto-migrate database tables
	if err := db.AutoMigrate(
		&entity.Role{},
		&entity.Permission{},
		&entity.User{},
		&entity.RefreshToken{},
		&entity.TenorLimit{},

		&entity.Consumer{},
		&entity.Transaction{},
		&entity.LimitMutation{},
	); err != nil {
		logger.SystemLogger.Fatal().Err(err).Msg("Failed to migrate database")
	}
	logger.SystemLogger.Info().Msg("Database migration completed successfully")

	database.SeedRBAC(app.DB)
	database.SeedUser(app.DB)
	database.SeedConsumerLimit(app.DB)
	database.SeedConsumer(app.DB)

}

// setupRouter initializes all dependencies and configures routes
func (app *Application) setupRouter() {

	userRepo := repository.NewUserRepository(app.DB)
	authService := services.NewAuthService(userRepo)
	refreshTokenRepo := repository.NewRefreshTokenRepository(app.DB)
	jwtService := services.NewJWTService(app.Config.JWT.Secret, app.Config.JWT.ExpiryHours, refreshTokenRepo)
	authHandler := handler.NewAuthHandler(authService, jwtService)

	limitRepo := repository.NewLimitRepository(app.DB)
	mutationRepo := repository.NewLimitMutationRepository(app.DB)
	limitService := services.NewLimitService(limitRepo, userRepo, mutationRepo, app.DB)
	limitHandler := handler.NewLimitHandler(limitService)
	userHandler := handler.NewUserHandler(userRepo)

	transactionRepo := repository.NewTransactionRepository(app.DB)
	transactionService := services.NewTransactionService(transactionRepo, limitRepo, mutationRepo, app.DB)
	transactionHandler := handler.NewTransactionHandler(transactionService)

	logService := services.NewLogService("logs") // Hardcoded log directory
	logHandler := handler.NewLogHandler(logService)

	appRouter := router.NewRouter(app.Config, authHandler, limitHandler, userHandler, transactionHandler, logHandler, userRepo)
	app.Router = appRouter.SetupRoutes()

	logger.SystemLogger.Info().Msg("Router configured successfully")
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
		logger.SystemLogger.Info().Msgf("🚀 Server running on port %s", app.Config.AppPort)
		if err := app.Server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logger.SystemLogger.Fatal().Err(err).Msg("❌ Server failed to start")
		}
	}()

	// Wait for interrupt signal for graceful shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	logger.SystemLogger.Info().Msg("Shutting down server...")

	// Graceful shutdown with 5 second timeout
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := app.Server.Shutdown(ctx); err != nil {
		logger.SystemLogger.Fatal().Err(err).Msg("Server forced to shutdown")
	}

	logger.SystemLogger.Info().Msg("Server exited successfully")
}
