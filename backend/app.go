package backend

import (
	"context"
	"dmarktodo/backend/models"
	"dmarktodo/backend/repository"
	"dmarktodo/backend/repository/inmemory"
	"log"
	"path/filepath"
)

// App struct
type App struct {
	ctx  context.Context
	repo repository.Repository
}

// NewApp creates a new App application struct
func NewApp() *App {
	repo, err := inmemory.NewRepository(filepath.Join("data", "tasks.txt"))
	if err != nil {
		log.Fatalf("Failed to create repository: %v", err)
	}
	return &App{
		ctx:  nil,
		repo: repo,
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
func (a *App) AddTask(title string) models.Task {
	return a.repo.AddTask(title)
}

// DeleteTask deletes a task by ID and returns true if successful
func (a *App) DeleteTask(id int) bool {
	return a.repo.DeleteTask(id)
}
