package main

import (
	"fmt"
	"log"

	"github.com/hadi-projects/xyz-finance-go/config"
	"github.com/hadi-projects/xyz-finance-go/pkg/database"
)

func main() {

	// 1. load config
	cfg, err := config.NewConfig()
	if err != nil {
		// Fail Fast: Aplikasi tidak boleh jalan jika config error
		log.Fatalf("Failed to load configuration: %v", err)
	}

	// 2. Init Database
	_, err = database.NewMySQLConnection(cfg)
	if err != nil {
		log.Fatalf("❌ Failed to connect to database: %v", err)
	}
	fmt.Println("✅ Database connected successfully")
}
