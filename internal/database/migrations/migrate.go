package migrations

import (
	"fmt"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/kurnosovmak/test/internal/config"
)

type Migrator struct {
	migrate *migrate.Migrate
}

func NewMigrator(cfg *config.Database, migrationsPath string) (*Migrator, error) {
	m, err := migrate.New(
		migrationsPath,
		cfg.GetDSN(),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create migrator: %v", err)
	}

	return &Migrator{migrate: m}, nil
}

func (m *Migrator) Up() error {
	if err := m.migrate.Up(); err != nil && err != migrate.ErrNoChange {
		return fmt.Errorf("failed to apply migrations: %v", err)
	}
	return nil
}

func (m *Migrator) Down() error {
	if err := m.migrate.Down(); err != nil && err != migrate.ErrNoChange {
		return fmt.Errorf("failed to rollback migrations: %v", err)
	}
	return nil
}

func (m *Migrator) Version() (uint, bool, error) {
	return m.migrate.Version()
}
