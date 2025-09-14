package errors

import "errors"

var (
	ErrDatabaseConnection = errors.New("database connection failed")
	ErrDatabaseMigration  = errors.New("database migration failed")
	ErrTaskNotFound      = errors.New("task not found")
	ErrInvalidTitle      = errors.New("title cannot be blank")
	ErrTaskValidation    = errors.New("task validation failed")
)