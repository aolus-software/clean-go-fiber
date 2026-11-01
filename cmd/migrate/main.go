package main

import (
	"log"
	"os"

	"github.com/aolus-software/clean-go-fiber/internal/config"
	"github.com/aolus-software/clean-go-fiber/internal/database/postgres"
)

func main() {
	config.LoadDatabaseOnly()

	if _, err := postgres.Connect(); err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	direction := "up"
	if len(os.Args) > 1 {
		direction = os.Args[1]
	}

	switch direction {
	case "up":
		if err := postgres.Migrate(); err != nil {
			log.Fatalf("Migration up failed: %v", err)
		}
		log.Println("Migrations applied")
	case "down":
		if err := postgres.MigrateDown(); err != nil {
			log.Fatalf("Migration down failed: %v", err)
		}
		log.Println("Migrations rolled back")
	default:
		log.Fatalf("Unknown direction %q: use up or down", direction)
	}
}
