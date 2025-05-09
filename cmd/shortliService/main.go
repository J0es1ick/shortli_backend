package main

import (
	"fmt"
	"log"

	"github.com/J0es1ick/shortli/internal/config"
	"github.com/J0es1ick/shortli/internal/database"
)

func main() {
	cfg, err := config.InitConfig()
	if err != nil {
		log.Fatalf("Config initialization error: %v", err)
	}

	fmt.Printf("Server port: %s\n", cfg.ServerPort)
	fmt.Printf("DB host: %s\n", cfg.Database.Host)

	db, err := database.DBInit(cfg)
	if err != nil {
		log.Fatalf("Failed to initialize database: %v", err)
	}
	fmt.Println("Connection successful")
	defer db.Close()
}