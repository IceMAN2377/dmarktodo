package postgres

import (
	"dmarktodo/backend/models"
	"dmarktodo/backend/repository"
	"github.com/jmoiron/sqlx"
	"log"
	"time"
)

type postgres struct {
	db *sqlx.DB
}

func NewRepository(db *sqlx.DB) repository.Repository {
	return &postgres{
		db: db,
	}
}

func (p *postgres) GetTasks() []models.Task {
	var tasks []models.Task
	err := p.db.Select(&tasks, "SELECT * FROM tasks ORDER BY id")
	if err != nil {
		// In a real application, you would handle this error properly
		// For now, we'll just log it and return an empty slice
		log.Printf("Error fetching tasks: %v", err)
		return []models.Task{}
	}
	return tasks
}

func (p *postgres) AddTask(title string) models.Task {
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
		log.Printf("Error adding task: %v", err)
		// Return an empty task in case of error
		return models.Task{}
	}

	return task
}

func (p *postgres) DeleteTask(id int) bool {
	// Delete the task with the given ID
	query := "DELETE FROM tasks WHERE id = $1"
	result, err := p.db.Exec(query, id)
	if err != nil {
		log.Printf("Error deleting task: %v", err)
		return false
	}

	// Check if any rows were affected (task was found and deleted)
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		log.Printf("Error getting rows affected: %v", err)
		return false
	}

	// Return true if at least one row was affected
	return rowsAffected > 0
}