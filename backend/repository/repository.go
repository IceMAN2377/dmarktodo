package repository

import (
	"dmarktodo/backend/models"
	"time"
)

type Repository interface {
	GetTasks(sortByCreatedDesc bool) []models.Task
	AddTask(task models.Task) (models.Task, error)
	DeleteTask(id int) bool
	UpdateTask(id int, status models.Status, completedAt *time.Time) (models.Task, error)
}
