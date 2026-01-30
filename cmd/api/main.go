package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/hadi-projects/xyz-finance-go/config"
	"github.com/hadi-projects/xyz-finance-go/internal/entity"
	"github.com/hadi-projects/xyz-finance-go/pkg/database"
)

func main() {

	// 1. load config
	cfg, err := config.NewConfig()
	if err != nil {
		log.Fatalf("Failed to load configuration: %v", err)
	}

	// 2. Init Database
	db, err := database.NewMySQLConnection(cfg)
	if err != nil {
		log.Fatalf("❌ Failed to connect to database: %v", err)
	}
	fmt.Println("✅ Database connected successfully")

	// 3. auto migrate
	if err := db.AutoMigrate(
		&entity.Role{},
		&entity.Permission{},
		&entity.User{},
	); err != nil {
		log.Fatalf("❌ Failed to migrate database: %v", err)
	}

	// 4. Init Router
	r := gin.New()
	r.Use(gin.Recovery())
	r.Use(gin.Logger())

	api := r.Group("/api/v1")
	{
		api.GET("/health", func(c *gin.Context) {
			// Cek Ping DB
			sqlDB, _ := db.DB()
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

	// 4. Run Server
	serverAddr := fmt.Sprintf(":%s", cfg.AppPort)
	fmt.Printf("🚀 Server running on port %s\n", cfg.AppPort)
	if err := r.Run(serverAddr); err != nil {
		log.Fatalf("❌ Server failed to start: %v", err)
	}
}
