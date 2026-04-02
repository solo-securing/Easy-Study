package main

import (
	"errors"
	"log"
	"os"

	"api/internal/database"
)

func main() {
	if len(os.Args) < 2 {
		log.Fatal("usage: go run ./cmd/migrate <up|down>")
	}

	cfg := database.MigrationConfig{
		SourceURL: os.Getenv("MIGRATIONS_SOURCE"),
		Host:      os.Getenv("BLUEPRINT_DB_HOST"),
		Port:      os.Getenv("BLUEPRINT_DB_PORT"),
		Database:  os.Getenv("BLUEPRINT_DB_DATABASE"),
		Username:  os.Getenv("BLUEPRINT_DB_USERNAME"),
		Password:  os.Getenv("BLUEPRINT_DB_PASSWORD"),
		SSLMode:   os.Getenv("BLUEPRINT_DB_SSLMODE"),
	}
	if cfg.SourceURL == "" {
		cfg.SourceURL = "file://migrations"
	}

	switch os.Args[1] {
	case "up":
		if err := database.RunMigrations(cfg); err != nil {
			log.Fatal(err)
		}
	case "down":
		m, err := database.NewMigrator(cfg)
		if err != nil {
			log.Fatal(err)
		}
		if err := m.Down(); err != nil && !errors.Is(err, database.ErrNoChange) {
			log.Fatal(err)
		}
		if err := m.Close(); err != nil {
			log.Fatal(err)
		}
	default:
		log.Fatal("unknown migrate command: " + os.Args[1])
	}
}
