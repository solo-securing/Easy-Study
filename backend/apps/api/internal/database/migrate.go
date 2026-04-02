package database

import (
	"errors"
	"fmt"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
)

var ErrNoChange = migrate.ErrNoChange

type MigrationConfig struct {
	SourceURL string
	Host      string
	Port      string
	Database  string
	Username  string
	Password  string
	SSLMode   string
}

func (c MigrationConfig) dsn() string {
	sslMode := c.SSLMode
	if sslMode == "" {
		sslMode = "disable"
	}

	return fmt.Sprintf(
		"postgres://%s:%s@%s:%s/%s?sslmode=%s",
		c.Username,
		c.Password,
		c.Host,
		c.Port,
		c.Database,
		sslMode,
	)
}

type Migrator struct {
	m *migrate.Migrate
}

func NewMigrator(cfg MigrationConfig) (*Migrator, error) {
	m, err := migrate.New(cfg.SourceURL, cfg.dsn())
	if err != nil {
		return nil, err
	}
	return &Migrator{m: m}, nil
}

func (m *Migrator) Up() error {
	return m.m.Up()
}

func (m *Migrator) Down() error {
	return m.m.Down()
}

func (m *Migrator) Close() error {
	sourceErr, dbErr := m.m.Close()
	if sourceErr != nil || dbErr != nil {
		return errors.Join(sourceErr, dbErr)
	}
	return nil
}

func RunMigrations(cfg MigrationConfig) error {
	m, err := NewMigrator(cfg)
	if err != nil {
		return err
	}

	if err := m.Up(); err != nil && !errors.Is(err, ErrNoChange) {
		return errors.Join(err, m.Close())
	}

	return m.Close()
}
