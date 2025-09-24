package handlers

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/NeroQue/course-management-backend/pkg/task"
)

// Response structs for task endpoints
type TaskResponse struct {
	Data *task.Task `json:"data"`
}

type TaskCleanupResponse struct {
	Message string `json:"message"`
	Cleaned int    `json:"cleaned"`
}

// TaskHandler handles task status requests
type TaskHandler struct{}

// NewTaskHandler creates new task handler
func NewTaskHandler() *TaskHandler {
	return &TaskHandler{}
}

// GetTask handles GET /api/tasks?id={taskId} - checks task status
func (h *TaskHandler) GetTask(w http.ResponseWriter, r *http.Request) {
	// Extract task ID from request
	taskID := r.URL.Query().Get("id")
	if taskID == "" {
		http.Error(w, "Task ID is required", http.StatusBadRequest)
		return
	}

	// check if task exists
	t, exists := task.GetTask(taskID)
	if !exists {
		http.Error(w, "Task not found", http.StatusNotFound)
		return
	}

	response := TaskResponse{
		Data: t,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

// CleanupTasks handles POST /api/tasks/cleanup - manually cleans old tasks
func (h *TaskHandler) CleanupTasks(w http.ResponseWriter, r *http.Request) {
	// default to 24 hours if not specified
	ageStr := r.URL.Query().Get("age")
	age := 24 * time.Hour

	if ageStr != "" {
		var err error
		age, err = time.ParseDuration(ageStr)
		if err != nil {
			http.Error(w, "Invalid duration format", http.StatusBadRequest)
			return
		}
	}

	// trigger cleanup
	cleaned := task.CleanupOldTasks(age)

	response := TaskCleanupResponse{
		Message: "Cleanup completed",
		Cleaned: cleaned,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}
