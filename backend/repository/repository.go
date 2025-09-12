package repository

import "dmarktodo/backend/models"

type Repository interface {
	GetTasks() []models.Task
	AddTask(title string) models.Task
	DeleteTask(id int) bool
}
