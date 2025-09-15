package service

import (
	apperrors "dmarktodo/backend/errors"
	"dmarktodo/backend/models"
	"dmarktodo/backend/repository"
	"log/slog"
	"strings"
	"time"
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

	// Create task with business logic defaults
	now := time.Now()
	task := models.Task{
		Title:     strings.TrimSpace(title),
		Status:    models.StatusActive,
		Priority:  priority,
		CreatedAt: now,
		UpdatedAt: now,
	}

	return s.repo.AddTask(task)
}

func (s *taskService) DeleteTask(id int) bool {
	return s.repo.DeleteTask(id)
}

func (s *taskService) ToggleStatus(id int, status models.Status) (models.Task, error) {
	// Business logic: determine new status and completedAt
	var newStatus models.Status
	var completedAt *time.Time

	if status == models.StatusActive {
		newStatus = models.StatusDone
		now := time.Now()
		completedAt = &now
	} else {
		newStatus = models.StatusActive
		completedAt = nil
	}

	return s.repo.UpdateTask(id, newStatus, completedAt)
}
