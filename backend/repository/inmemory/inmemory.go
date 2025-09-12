package inmemory

import (
	"bufio"
	"dmarktodo/backend/models"
	"dmarktodo/backend/repository"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"sync"
)

// InMemoryRepo implements the repository.Repository interface
type inMemoryRepo struct {
	tasks    []models.Task
	nextID   int
	filePath string
	mutex    sync.RWMutex
}

// NewRepository creates a new instance of InMemoryRepo
func NewRepository(filePath string) (repository.Repository, error) {
	// Create directory if it doesn't exist
	dir := filepath.Dir(filePath)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return nil, fmt.Errorf("failed to create directory: %w", err)
	}

	repo := &inMemoryRepo{
		tasks:    make([]models.Task, 0),
		nextID:   1,
		filePath: filePath,
	}

	// Load existing tasks from file if it exists
	if err := repo.loadFromFile(); err != nil {
		return nil, err
	}

	return repo, nil
}

// GetTasks returns all tasks
func (r *inMemoryRepo) GetTasks() []models.Task {
	r.mutex.RLock()
	defer r.mutex.RUnlock()

	// Return a copy of the tasks slice to prevent modification from outside
	tasksCopy := make([]models.Task, len(r.tasks))
	copy(tasksCopy, r.tasks)
	return tasksCopy
}

// AddTask adds a new task and returns it
func (r *inMemoryRepo) AddTask(title string) models.Task {
	r.mutex.Lock()
	defer r.mutex.Unlock()

	task := models.Task{
		ID:    r.nextID,
		Title: title,
	}

	r.tasks = append(r.tasks, task)
	r.nextID++

	// Save to file after adding a task
	_ = r.saveToFile() // Ignoring error for simplicity, but in production you should handle it

	return task
}

// loadFromFile loads tasks from the file
func (r *inMemoryRepo) loadFromFile() error {
	// Check if file exists
	if _, err := os.Stat(r.filePath); os.IsNotExist(err) {
		// File doesn't exist, which is fine for a new repository
		return nil
	}

	file, err := os.Open(r.filePath)
	if err != nil {
		return fmt.Errorf("failed to open file: %w", err)
	}
	defer file.Close()

	r.tasks = make([]models.Task, 0)
	maxID := 0

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		var task models.Task
		if err := json.Unmarshal([]byte(line), &task); err != nil {
			return fmt.Errorf("failed to unmarshal task: %w", err)
		}

		r.tasks = append(r.tasks, task)

		// Keep track of the highest ID
		if task.ID > maxID {
			maxID = task.ID
		}
	}

	if err := scanner.Err(); err != nil {
		return fmt.Errorf("error reading file: %w", err)
	}

	// Set nextID to be one more than the highest ID found
	r.nextID = maxID + 1

	return nil
}

// DeleteTask deletes a task by ID and returns true if successful
func (r *inMemoryRepo) DeleteTask(id int) bool {
	r.mutex.Lock()
	defer r.mutex.Unlock()

	for i, task := range r.tasks {
		if task.ID == id {
			// Remove the task from the slice
			r.tasks = append(r.tasks[:i], r.tasks[i+1:]...)
			
			// Save to file after deleting a task
			_ = r.saveToFile() // Ignoring error for simplicity
			
			return true
		}
	}
	
	return false
}

// saveToFile saves tasks to the file
func (r *inMemoryRepo) saveToFile() error {
	file, err := os.Create(r.filePath)
	if err != nil {
		return fmt.Errorf("failed to create file: %w", err)
	}
	defer file.Close()

	writer := bufio.NewWriter(file)
	for _, task := range r.tasks {
		taskJSON, err := json.Marshal(task)
		if err != nil {
			return fmt.Errorf("failed to marshal task: %w", err)
		}

		if _, err := writer.WriteString(string(taskJSON) + "\n"); err != nil {
			return fmt.Errorf("failed to write to file: %w", err)
		}
	}

	if err := writer.Flush(); err != nil {
		return fmt.Errorf("failed to flush writer: %w", err)
	}

	return nil
}
