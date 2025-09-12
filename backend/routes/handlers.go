package routes

import (
	"dmarktodo/backend/models"
	"dmarktodo/backend/repository"
)

type RecordHandler struct {
	repo repository.Repository
}

func NewRecordHandler(repo repository.Repository) *RecordHandler {
	return &RecordHandler{
		repo: repo,
	}
}

func (r *RecordHandler) GetTasks() []models.Task {
	return r.repo.GetTasks()
}

func (r *RecordHandler) AddTask(title string) models.Task {
	return r.repo.AddTask(title)
}
