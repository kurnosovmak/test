package main

import (
	"flag"
	"log"

	"github.com/golang-migrate/migrate/v4"
	"github.com/kurnosovmak/test/internal/config"
	"github.com/kurnosovmak/test/internal/database/migrations"
)

const (
	migrationsPath = "file://migrations"
)

var (
	envPath = flag.String("config", ".yaml", "Path to the environment file")
	action  = flag.String("action", "up", "Migration action (up/down/version)")
)

func init() {
	flag.Parse()
}

func main() {
	// Загружаем конфигурацию
	cfg, err := config.LoadConfig(*envPath)
	if err != nil {
		log.Fatalf("Error loading config: %v", err)
	}

	m, err := migrations.NewMigrator(&cfg.Database, migrationsPath)
	if err != nil {
		log.Fatalf("Error creating migrator: %v", err)
	}

	// Выполняем действие в зависимости от флага
	switch *action {
	case "up":
		if err := m.Up(); err != nil && err != migrate.ErrNoChange {
			log.Fatalf("Error applying migrations: %v", err)
		}
		log.Println("Successfully applied all migrations")

	case "down":
		if err := m.Down(); err != nil && err != migrate.ErrNoChange {
			log.Fatalf("Error rolling back migrations: %v", err)
		}
		log.Println("Successfully rolled back all migrations")

	case "version":
		version, dirty, err := m.Version()
		if err != nil {
			log.Fatalf("Error getting migration version: %v", err)
		}
		log.Printf("Current migration version: %d, Dirty: %v", version, dirty)

	default:
		log.Fatalf("Unknown action: %s", action)
	}
}
