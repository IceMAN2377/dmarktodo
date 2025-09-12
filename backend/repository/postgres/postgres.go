package postgres

import (
	"dmarktodo/backend/models"
	"github.com/jmoiron/sqlx"
)

type postgres struct {
	db *sqlx.DB
}

func NewRepository(db *sqlx.DB) *postgres {
	return &postgres{
		db: db,
	}
}

func (p *postgres) GetTasks() []models.Task {
	return nil
}

func (p *postgres) CreateTask(task models.Task) error {
	return nil
}
