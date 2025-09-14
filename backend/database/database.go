package database

import (
	"database/sql"
	"dmarktodo/backend/config"
	apperrors "dmarktodo/backend/errors"
	"errors"
	"fmt"
	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"net/url"
)

const (
	migrationsPath = "db/migrations"
	dbDriverName   = "postgres"
)

type DB struct {
	*sqlx.DB
}

// NewConnection creates a new database connection and runs migrations if enabled
func NewConnection(config *config.Config) (*DB, error) {
	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=%s",
		config.PostgresHost, config.PostgresPort, config.PostgresUser, config.PostgresPassword, config.PostgresDb, config.PostgresSslMode)

	db, err := sql.Open(dbDriverName, psqlInfo)
	if err != nil {
		return nil, fmt.Errorf("%w: %v", apperrors.ErrDatabaseConnection, err)
	}

	psql := sqlx.NewDb(db, dbDriverName)

	// Run migrations if enabled
	if config.PostgresMigrate {
		if err := runMigrations(config); err != nil {
			return nil, fmt.Errorf("%w: %v", apperrors.ErrDatabaseMigration, err)
		}
	}

	return &DB{psql}, nil
}

func runMigrations(config *config.Config) error {
	dbURI := fmt.Sprintf("postgres://%s:%s@%s:%d/%s?sslmode=%s",
		config.PostgresUser, url.QueryEscape(config.PostgresPassword), config.PostgresHost, config.PostgresPort, config.PostgresDb, config.PostgresSslMode)

	m, err := migrate.New("file://"+migrationsPath, dbURI)
	if err != nil {
		return fmt.Errorf("failed to create migrate object: %w", err)
	}

	if err := m.Up(); err != nil && !errors.Is(err, migrate.ErrNoChange) {
		return fmt.Errorf("failed to apply migrations: %w", err)
	}

	return nil
}