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
	); err != nil {
		log.Fatalf("Failed to migrate database: %v", err)
	}
	log.Println("Database migration completed successfully")
}

// run starts the HTTP server and handles graceful shutdown
func (app *Application) run() {
	// 3. Init Router
	r := gin.New()
	r.Use(gin.Recovery())
	r.Use(gin.Logger())

	api := r.Group("/api/v1")
	{
		api.GET("/health", func(c *gin.Context) {
			// Cek Ping DB
			sqlDB, _ := app.DB.DB()
			status := "UP"
			if err := sqlDB.Ping(); err != nil {
				status = "DOWN (DB Error)"
			}

			c.JSON(http.StatusOK, gin.H{
				"status":  status,
				"app":     "XYZ Multifinance",
				"version": "1.0.0",
			})
		})
	}
	app.Router = r

	// Start server in goroutine
	go func() {
		serverAddr := fmt.Sprintf(":%s", app.Config.AppPort)
		fmt.Printf("🚀 Server running on port %s\n", app.Config.AppPort)
		if err := app.Router.Run(serverAddr); err != nil {
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
