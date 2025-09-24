package task

import (
	"sync"
	"time"

	"github.com/google/uuid"
)

// Status shows what state a task is in
type Status string

const (
	StatusPending    Status = "pending"    // waiting to start
	StatusProcessing Status = "processing" // currently running
	StatusCompleted  Status = "completed"  // finished successfully
	StatusFailed     Status = "failed"     // something went wrong
)

// Task represents a background job that might take a while
type Task struct {
	ID           string      `json:"id"`
	Type         string      `json:"type"`                    // what kind of task
	Status       Status      `json:"status"`                  // current state
	Progress     float32     `json:"progress"`                // 0-100 percent done
	CreatedAt    time.Time   `json:"created_at"`              // when it started
	StartedAt    time.Time   `json:"started_at,omitempty"`    // when processing began
	CompletedAt  time.Time   `json:"completed_at,omitempty"`  // when it finished
	Message      string      `json:"message,omitempty"`       // status updates
	ErrorMessage string      `json:"error_message,omitempty"` // what went wrong
	Result       interface{} `json:"result,omitempty"`        // final results
}

// TaskManager keeps track of all running tasks
type TaskManager struct {
	tasks map[string]*Task
	mu    sync.RWMutex // for thread safety
}

// global task manager - another singleton but whatever
var manager *TaskManager

// Initialize sets up the task manager
func Initialize() {
	manager = &TaskManager{
		tasks: make(map[string]*Task),
	}
}

// CreateTask makes a new task and returns its ID
func CreateTask(taskType string) string {
	if manager == nil {
		Initialize()
	}

	taskID := uuid.New().String()
	task := &Task{
		ID:        taskID,
		Type:      taskType,
		Status:    StatusPending,
		Progress:  0,
		CreatedAt: time.Now(),
	}

	manager.mu.Lock()
	manager.tasks[taskID] = task
	manager.mu.Unlock()

	return taskID
}

// GetTask retrieves task info by ID
func GetTask(taskID string) (*Task, bool) {
	if manager == nil {
		return nil, false
	}

	manager.mu.RLock()
	defer manager.mu.RUnlock()

	task, exists := manager.tasks[taskID]
	return task, exists
}

// UpdateTaskStatus changes the task status
func UpdateTaskStatus(taskID string, status Status) {
	if manager == nil {
		return
	}

	manager.mu.Lock()
	defer manager.mu.Unlock()

	task, exists := manager.tasks[taskID]
	if !exists {
		return
	}

	task.Status = status
	if status == StatusProcessing && task.StartedAt.IsZero() {
		task.StartedAt = time.Now()
	}
	if status == StatusCompleted || status == StatusFailed {
		task.CompletedAt = time.Now()
	}
}

// UpdateTaskProgress updates how much of the task is done
func UpdateTaskProgress(taskID string, progress float32, message string) {
	if manager == nil {
		return
	}

	manager.mu.Lock()
	defer manager.mu.Unlock()

	task, exists := manager.tasks[taskID]
	if !exists {
		return
	}

	task.Progress = progress
	task.Message = message
}

// SetTaskMessage updates the status message
func SetTaskMessage(taskID string, message string) {
	if manager == nil {
		return
	}

	manager.mu.Lock()
	defer manager.mu.Unlock()

	task, exists := manager.tasks[taskID]
	if !exists {
		return
	}

	task.Message = message
}

// SetTaskError marks task as failed with error message
func SetTaskError(taskID string, errorMessage string) {
	if manager == nil {
		return
	}

	manager.mu.Lock()
	defer manager.mu.Unlock()

	task, exists := manager.tasks[taskID]
	if !exists {
		return
	}

	task.Status = StatusFailed
	task.ErrorMessage = errorMessage
	task.CompletedAt = time.Now()
}

// CompleteTask marks task as done with optional result data
func CompleteTask(taskID string, result interface{}) {
	if manager == nil {
		return
	}

	manager.mu.Lock()
	defer manager.mu.Unlock()

	task, exists := manager.tasks[taskID]
	if !exists {
		return
	}

	task.Status = StatusCompleted
	task.Progress = 100
	task.Result = result
	task.CompletedAt = time.Now()
}

// CleanupOldTasks removes completed tasks older than the specified age
func CleanupOldTasks(maxAge time.Duration) int {
	if manager == nil {
		return 0
	}

	manager.mu.Lock()
	defer manager.mu.Unlock()

	cutoff := time.Now().Add(-maxAge)
	cleaned := 0

	for taskID, task := range manager.tasks {
		// only clean up completed or failed tasks
		if (task.Status == StatusCompleted || task.Status == StatusFailed) &&
			!task.CompletedAt.IsZero() && task.CompletedAt.Before(cutoff) {
			delete(manager.tasks, taskID)
			cleaned++
		}
	}

	return cleaned
}

// CleanupRoutine runs cleanup automatically on a schedule
func CleanupRoutine(interval, maxAge time.Duration) {
	ticker := time.NewTicker(interval)
	defer ticker.Stop()

	for range ticker.C {
		cleaned := CleanupOldTasks(maxAge)
		if cleaned > 0 {
			// maybe log this but don't spam the logs
		}
	}
}
