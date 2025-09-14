package service

import (
	apperrors "dmarktodo/backend/errors"
	"dmarktodo/backend/models"
	"dmarktodo/backend/repository"
	"log/slog"
	"strings"
)

type TaskService interface {
	GetTasks(sortByCreatedDesc bool) []models.Task
	AddTask(title string, priority models.Priority) (models.Task, error)
	DeleteTask(id int) bool
	ToggleStatus(id int, status models.Status) (models.Task, error)
}

type taskService struct {
	logger *slog.Logger
	repo   repository.Repository
}

func NewTaskService(logger *slog.Logger, repo repository.Repository) TaskService {
	return &taskService{
		logger: logger,
		repo:   repo,
	}
}

func (s *taskService) GetTasks(sortByCreatedDesc bool) []models.Task {
	return s.repo.GetTasks(sortByCreatedDesc)
}

func (s *taskService) AddTask(title string, priority models.Priority) (models.Task, error) {
	if strings.TrimSpace(title) == "" {
		s.logger.Error("title cannot be blank")
		return models.Task{}, apperrors.ErrInvalidTitle
	}
	return s.repo.AddTask(title, priority)
}

func (s *taskService) DeleteTask(id int) bool {
	return s.repo.DeleteTask(id)
}

func (s *taskService) ToggleStatus(id int, status models.Status) (models.Task, error) {
	return s.repo.ToggleStatus(id, status)
}
