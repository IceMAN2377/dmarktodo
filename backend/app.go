package backend

import (
	"context"
	"dmarktodo/backend/config"
	"dmarktodo/backend/database"
	"dmarktodo/backend/models"
	"dmarktodo/backend/repository/postgres"
	"dmarktodo/backend/service"
	"fmt"
	"log/slog"
	"os"
)

// App struct
type App struct {
	ctx     context.Context
	logger  *slog.Logger
	service service.TaskService
}

// NewApp creates a new App application struct
func NewApp(config *config.Config) (*App, error) {
	logger := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}))

	db, err := database.NewConnection(config)
	if err != nil {
		return nil, fmt.Errorf("failed to initialize app: %w", err)
	}

	repo := postgres.NewRepository(logger, db.DB)
	taskService := service.NewTaskService(logger, repo)

	return &App{
		ctx:     nil,
		logger:  logger,
		service: taskService,
	}, nil
}

// startup is called when the app starts. The context is saved
// so we can call the runtime methods
func (a *App) Startup(ctx context.Context) {
	a.ctx = ctx
}

func (a *App) GetTasks(sortByCreatedDesc bool) []models.Task {
	return a.service.GetTasks(sortByCreatedDesc)
}

// AddTask adds a new task and returns it
func (a *App) AddTask(title string, priority models.Priority) (models.Task, error) {
	return a.service.AddTask(title, priority)
}

// DeleteTask deletes a task by ID and returns true if successful
func (a *App) DeleteTask(id int) bool {
	return a.service.DeleteTask(id)
}

func (a *App) ToggleStatus(id int, status models.Status) (models.Task, error) {
	return a.service.ToggleStatus(id, status)
}
