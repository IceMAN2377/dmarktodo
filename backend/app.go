package backend

import (
	"context"
)

// App struct
type App struct {
	ctx    context.Context
	tasks  []Task
	nextID int
}

// NewApp creates a new App application struct
func NewApp() *App {
	return &App{
		tasks:  []Task{},
		nextID: 1,
	}
}

// startup is called when the app starts. The context is saved
// so we can call the runtime methods
func (a *App) Startup(ctx context.Context) {
	a.ctx = ctx
}

type Task struct {
	ID    int    `json:"id"`
	Title string `json:"title"`
}

func (a *App) GetTasks() []Task {
	return a.tasks
}

func (a *App) AddTask(title string) Task {
	task := Task{
		ID:    a.nextID,
		Title: title,
	}

	a.tasks = append(a.tasks, task)
	a.nextID++

	return task
}
