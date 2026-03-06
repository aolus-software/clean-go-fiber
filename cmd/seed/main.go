package main

import (
	"log"

	"github.com/aolus-software/clean-go-fiber/internal/config"
	"github.com/aolus-software/clean-go-fiber/internal/database/postgres"
	"github.com/aolus-software/clean-go-fiber/internal/database/postgres/seed"
)

func main() {
	config.Load()

	if _, err := postgres.Connect(); err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	if err := seed.Run(); err != nil {
		log.Fatalf("Failed to run seeders: %v", err)
	}

	log.Println("Seeders applied successfully")
}
