package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/J0es1ick/shortli/internal/app/routes"
	"github.com/J0es1ick/shortli/internal/config"
	"github.com/J0es1ick/shortli/internal/database"
	"github.com/J0es1ick/shortli/internal/repository"
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

	urlRepo := repository.NewUrlRepository(db.DB)
	handler := routes.SetupRoutes(cfg, urlRepo)

	server := &http.Server{
		Addr:    ":" + cfg.ServerPort,
		Handler: handler,
	}

	log.Printf("Server starting on port %s", cfg.ServerPort)
	if err := server.ListenAndServe(); err != nil {
		log.Fatalf("Server failed: %v", err)
	}
}