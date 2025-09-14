package backend

import (
	"context"
	"database/sql"
	"dmarktodo/backend/config"
	"dmarktodo/backend/models"
	"dmarktodo/backend/repository"
	"dmarktodo/backend/repository/postgres"
	"errors"
	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"log/slog"
	"os"
	"strings"

	"fmt"
	"github.com/jmoiron/sqlx"
	"net/url"

	_ "github.com/lib/pq"
)

const (
	migrationsPath = "db/migrations"
	dbDriverName   = "postgres"
)

// App struct
type App struct {
	ctx    context.Context
	logger *slog.Logger
	repo   repository.Repository
}

// NewApp creates a new App application struct
func NewApp(config *config.Config) *App {

	logger := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}))

	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=%s",
		config.PostgresHost, config.PostgresPort, config.PostgresUser, config.PostgresPassword, config.PostgresDb, config.PostgresSslMode)

	db, err := sql.Open(dbDriverName, psqlInfo)
	if err != nil {
		panic("failed to connect to db: " + err.Error())
	}

	psql := sqlx.NewDb(db, dbDriverName)

	dbURI := fmt.Sprintf("postgres://%s:%s@%s:%d/%s?sslmode=%s",
		config.PostgresUser, url.QueryEscape(config.PostgresPassword), config.PostgresHost, config.PostgresPort, config.PostgresDb, config.PostgresSslMode)
	if config.PostgresMigrate {
		if m, err := migrate.New("file://"+migrationsPath, dbURI); err != nil {
			panic("failed to create migrate object: " + err.Error())
		} else if err = m.Up(); err != nil && !errors.Is(err, migrate.ErrNoChange) {
			panic("failed to apply migrations: " + err.Error())
		}
	}

	repo := postgres.NewRepository(logger, psql)
	return &App{
		ctx:    nil,
		logger: logger,
		repo:   repo,
	}
}

// startup is called when the app starts. The context is saved
// so we can call the runtime methods
func (a *App) Startup(ctx context.Context) {
	a.ctx = ctx
}

func (a *App) GetTasks() []models.Task {
	return a.repo.GetTasks()
}

// AddTask adds a new task and returns it
func (a *App) AddTask(title string) (models.Task, error) {

	if strings.TrimSpace(title) == "" {
		a.logger.Error("title cannot be blank")
		return models.Task{}, errors.New("title cannot be blank")
	}
	return a.repo.AddTask(title)
}

// DeleteTask deletes a task by ID and returns true if successful
func (a *App) DeleteTask(id int) bool {
	return a.repo.DeleteTask(id)
}

func (a *App) ToggleStatus(id int, status models.Status) (models.Task, error) {
	return a.repo.ToggleStatus(id, status)
}
