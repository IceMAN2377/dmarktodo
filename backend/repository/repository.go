package repository

import "dmarktodo/backend/models"

type Repository interface {
	GetTasks(sortByCreatedDesc bool) []models.Task
	AddTask(title string, priority models.Priority) (models.Task, error)
	DeleteTask(id int) bool
	ToggleStatus(id int, status models.Status) (models.Task, error)
}
