package handlers

import (
	"log"
	"net/http"
	"time"

	"github.com/NeroQue/course-management-backend/pkg/task"
)

// TaskHandler handles task status requests
type TaskHandler struct{}

// NewTaskHandler creates new task handler
func NewTaskHandler() *TaskHandler {
	return &TaskHandler{}
}

// GetTask handles GET /api/tasks?id={taskId} - checks task status
// @Summary Get task status
// @Description Get the current status and result of a background task
// @Tags tasks
// @Accept json
// @Produce json
// @Param id query string true "Task ID"
// @Success 200 {object} map[string]interface{} "Task status and information"
// @Failure 400 {object} map[string]interface{} "Missing or invalid task ID"
// @Failure 404 {object} map[string]interface{} "Task not found"
// @Router /tasks [get]
func (h *TaskHandler) GetTask(w http.ResponseWriter, r *http.Request) {
	log.Printf("Task status requested from IP: %s", r.RemoteAddr)

	// Extract task ID from request
	taskID := r.URL.Query().Get("id")
	if taskID == "" {
		SendErrorResponse(w, "Task ID is required", http.StatusBadRequest,
			"Task status request without task ID", nil)
		return
	}

	log.Printf("Looking up task: %s", taskID)

	// check if task exists
	t, exists := task.GetTask(taskID)
	if !exists {
		SendErrorResponse(w, "Task not found", http.StatusNotFound,
			"Requested task does not exist: "+taskID, nil)
		return
	}

	SendSuccessResponse(w, "Task status retrieved", t,
		"Task status retrieved for: "+taskID)
}

// CleanupTasks handles POST /api/tasks/cleanup - manually cleans old tasks
// @Summary Cleanup old tasks
// @Description Manually trigger cleanup of old completed tasks
// @Tags tasks
// @Accept json
// @Produce json
// @Param age query string false "Age duration (e.g., '24h', '1h30m'). Defaults to 24h"
// @Success 200 {object} map[string]interface{} "Number of tasks cleaned up"
// @Failure 400 {object} map[string]interface{} "Invalid duration format"
// @Router /tasks/cleanup [post]
func (h *TaskHandler) CleanupTasks(w http.ResponseWriter, r *http.Request) {
	log.Printf("Task cleanup requested from IP: %s", r.RemoteAddr)

	// default to 24 hours if not specified
	ageStr := r.URL.Query().Get("age")
	age := 24 * time.Hour

	if ageStr != "" {
		var err error
		age, err = time.ParseDuration(ageStr)
		if err != nil {
			SendErrorResponse(w, "Invalid duration format", http.StatusBadRequest,
				"Invalid age duration in task cleanup: "+ageStr, err)
			return
		}
	}

	log.Printf("Starting task cleanup for tasks older than: %v", age)

	// trigger cleanup
	cleaned := task.CleanupOldTasks(age)

	responseData := map[string]interface{}{
		"cleaned": cleaned,
		"age":     age.String(),
	}

	SendSuccessResponse(w, "Cleanup completed", responseData,
		"Task cleanup completed - cleaned "+string(rune(cleaned))+" tasks")
}
