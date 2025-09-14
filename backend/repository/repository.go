package repository

import "dmarktodo/backend/models"

type Repository interface {
	GetTasks() []models.Task
	AddTask(title string) (models.Task, error)
	DeleteTask(id int) bool
	ToggleStatus(id int, status models.Status) (models.Task, error)
}
