package postgres

import (
	"database/sql"
	"dmarktodo/backend/models"
	"dmarktodo/backend/repository"
	"errors"
	"github.com/jmoiron/sqlx"
	"log/slog"
	"time"
)

type postgres struct {
	logger *slog.Logger
	db     *sqlx.DB
}

func NewRepository(logger *slog.Logger, db *sqlx.DB) repository.Repository {
	return &postgres{
		logger: logger,
		db:     db,
	}
}

func (p *postgres) GetTasks() []models.Task {
	var tasks []models.Task
	err := p.db.Select(&tasks, "SELECT * FROM tasks ORDER BY id")
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			p.logger.Info("No tasks found - this is normal for a new database")
			return []models.Task{}
		}
	}
	return tasks
}

func (p *postgres) AddTask(title string) (models.Task, error) {
	// Create a new task with default values
	now := time.Now()
	task := models.Task{
		Title:     title,
		Status:    models.StatusActive,
		Priority:  models.PriorityMedium,
		CreatedAt: now,
		UpdatedAt: now,
	}

	// Insert the task into the database and get the ID
	query := `
		INSERT INTO tasks (title, status, priority, created_at, updated_at) 
		VALUES ($1, $2, $3, $4, $5) 
		RETURNING id
	`
	err := p.db.QueryRow(
		query,
		task.Title,
		task.Status,
		task.Priority,
		task.CreatedAt,
		task.UpdatedAt,
	).Scan(&task.ID)

	if err != nil {
		// In a real application, you would handle this error properly
		p.logger.Error("Error adding task")
		// Return an empty task in case of error
		return models.Task{}, nil
	}

	return task, nil
}

func (p *postgres) DeleteTask(id int) bool {
	// Delete the task with the given ID
	query := "DELETE FROM tasks WHERE id = $1"
	result, err := p.db.Exec(query, id)
	if err != nil {
		p.logger.Error("Error deleting task")
		return false
	}

	// Check if any rows were affected (task was found and deleted)
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		p.logger.Error("Error getting rows affected")
		return false
	}

	// Return true if at least one row was affected
	return rowsAffected > 0
}

func (p *postgres) ToggleStatus(id int, status models.Status) (models.Task, error) {
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

	// Обновляем статус задачи
	query := `UPDATE tasks 
		SET status = $1, completed_at = $2, updated_at = NOW() 
		WHERE id = $3 
		RETURNING id, title, status, priority, completed_at, created_at, updated_at`

	var task models.Task
	err := p.db.QueryRow(query, newStatus, completedAt, id).Scan(
		&task.ID,
		&task.Title,
		&task.Status,
		&task.Priority,
		&task.CompletedAt,
		&task.CreatedAt,
		&task.UpdatedAt,
	)

	if err != nil {
		p.logger.Error("Error updating task status", "error", err, "id", id)
		return models.Task{}, err
	}

	return task, nil
}
