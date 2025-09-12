package models

import (
	"time"
)

// Статус задачи
type Status string

const (
	StatusActive Status = "active"
	StatusDone   Status = "done"
)

// Приоритет задачи
type Priority string

const (
	PriorityLow    Priority = "low"
	PriorityMedium Priority = "medium"
	PriorityHigh   Priority = "high"
)

type Task struct {
	ID          int        `db:"id" json:"id"`
	Title       string     `db:"title" json:"title"`
	Status      Status     `db:"status" json:"status"`
	Priority    Priority   `db:"priority" json:"priority"`
	CompletedAt *time.Time `db:"completed_at" json:"completed_at,omitempty"`

	CreatedAt time.Time `db:"created_at" json:"created_at"`
	UpdatedAt time.Time `db:"updated_at" json:"updated_at"`
}
