package handlers

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/NeroQue/course-management-backend/internal/models"
	"github.com/NeroQue/course-management-backend/internal/services"
	"github.com/NeroQue/course-management-backend/pkg/parser"
	"github.com/NeroQue/course-management-backend/pkg/session"
	"github.com/NeroQue/course-management-backend/pkg/task"
	"github.com/NeroQue/course-management-backend/pkg/util"
	"github.com/google/uuid"
)

// Response structs for course endpoints
type CourseListResponse struct {
	Message string           `json:"message"`
	Data    []*models.Course `json:"data"`
}

type CourseResponse struct {
	Message string         `json:"message"`
	Data    *models.Course `json:"data"`
}

type DirectoriesResponse struct {
	Message string            `json:"message"`
	Data    []parser.FileInfo `json:"data"`
}

type ScanCoursesResponse struct {
	Message string            `json:"message"`
	Count   int               `json:"count"`
	Data    []parser.FileInfo `json:"data"`
}

type BatchImportStartResponse struct {
	Message string `json:"message"`
	TaskID  string `json:"task_id"`
}

type CourseProgressResponse struct {
	Message string                 `json:"message"`
	Data    *models.CourseProgress `json:"data"`
}

type ModuleProgressResponse struct {
	Message string                 `json:"message"`
	Data    *models.ModuleProgress `json:"data"`
}

type ProgressUpdateResponse struct {
	Message string `json:"message"`
}

type ContentCompleteResponse struct {
	Message string `json:"message"`
}

type UserProgressSummaryResponse struct {
	Message string                  `json:"message"`
	Data    *models.ProgressSummary `json:"data"`
}

// request/response structs for batch import
type BatchImportRequest struct {
	Courses []models.CreateCourseInput `json:"courses"`
}

type BatchImportResponse struct {
	SuccessCount    int              `json:"success_count"`
	FailureCount    int              `json:"failure_count"`
	ImportedCourses []*models.Course `json:"imported_courses"`
	Errors          []string         `json:"errors,omitempty"`
}

// CourseHandler processes course-related HTTP requests
type CourseHandler struct {
	Service *services.CourseService // handles all course business logic
}

// NewCourseHandler creates handler with injected service
func NewCourseHandler(service *services.CourseService) *CourseHandler {
	return &CourseHandler{Service: service}
}

// List handles GET /api/courses - returns all courses
func (h *CourseHandler) List(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	// get courses from service layer
	courses, err := h.Service.ListCourses(r.Context())
	if err != nil {
		log.Printf("Error retrieving courses: %v", err)
		http.Error(w, "Failed to retrieve courses", http.StatusInternalServerError)
		return
	}

	response := CourseListResponse{
		Message: "Courses retrieved successfully",
		Data:    courses,
	}

	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(response); err != nil {
		log.Printf("Error encoding JSON response: %v", err)
		http.Error(w, "Failed to encode response", http.StatusInternalServerError)
		return
	}
}

// Create handles POST /api/courses - makes new course from directory
func (h *CourseHandler) Create(w http.ResponseWriter, r *http.Request) {
	var input models.CreateCourseInput
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// basic validation
	if input.Title == "" {
		http.Error(w, "Course title is required", http.StatusBadRequest)
		return
	}

	if input.RelativePath == "" {
		http.Error(w, "Relative path is required", http.StatusBadRequest)
		return
	}

	// need user logged in to create courses
	userID := session.GetCurrentUser()
	if userID == uuid.Nil {
		http.Error(w, "You must be logged in to create courses", http.StatusUnauthorized)
		return
	}

	if input.BasePath == "" {
		input.BasePath = util.GetCoursesDirectory()
	}

	directoryPath := filepath.Join(input.BasePath, input.RelativePath)
	ctx := r.Context()

	// let service handle the actual import
	course, err := h.Service.ImportCourse(ctx, directoryPath, userID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	response := CourseResponse{
		Message: "Course created successfully",
		Data:    course,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(response)
}

// ListDirectories handles GET /api/courses/directories - shows available dirs
func (h *CourseHandler) ListDirectories(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	directories, err := h.Service.Parser.ListCourseDirectories()
	if err != nil {
		log.Printf("Error listing directories: %v", err)
		http.Error(w, "Failed to list directories", http.StatusInternalServerError)
		return
	}

	response := DirectoriesResponse{
		Message: "Directories retrieved successfully",
		Data:    directories,
	}

	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(response); err != nil {
		log.Printf("Error encoding JSON response: %v", err)
		http.Error(w, "Failed to encode response", http.StatusInternalServerError)
		return
	}
}

// ScanNewCourses handles GET /api/courses/scan - finds dirs not imported yet
func (h *CourseHandler) ScanNewCourses(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	// compare filesystem with database to find new ones
	newDirectories, err := h.Service.ScanNewCourses(r.Context())
	if err != nil {
		log.Printf("Error scanning for new courses: %v", err)
		http.Error(w, "Failed to scan for new courses", http.StatusInternalServerError)
		return
	}

	response := ScanCoursesResponse{
		Message: "New course directories found",
		Count:   len(newDirectories),
		Data:    newDirectories,
	}

	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(response); err != nil {
		log.Printf("Error encoding JSON response: %v", err)
		http.Error(w, "Failed to encode response", http.StatusInternalServerError)
		return
	}
}

// BatchImport handles POST /api/courses/batch - imports multiple courses at once
func (h *CourseHandler) BatchImport(w http.ResponseWriter, r *http.Request) {
	var request BatchImportRequest
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	if len(request.Courses) == 0 {
		http.Error(w, "No courses provided for import", http.StatusBadRequest)
		return
	}

	userID := session.GetCurrentUser()
	if userID == uuid.Nil {
		http.Error(w, "You must be logged in to import courses", http.StatusUnauthorized)
		return
	}

	// create background task since this might take a while
	taskID := task.CreateTask("batch_import")

	// do the actual work in background
	go func() {
		task.UpdateTaskStatus(taskID, task.StatusProcessing)
		task.SetTaskMessage(taskID, "Starting import of "+strconv.Itoa(len(request.Courses))+" courses")

		// need new context since original request will be done
		ctx := context.Background()

		importedCourses, errs := h.Service.BatchImportCourses(ctx, request.Courses, userID)

		response := BatchImportResponse{
			SuccessCount:    len(importedCourses),
			FailureCount:    len(errs),
			ImportedCourses: importedCourses,
		}

		for _, err := range errs {
			response.Errors = append(response.Errors, err.Error())
		}

		// update task based on results
		if len(errs) > 0 && len(importedCourses) == 0 {
			task.SetTaskError(taskID, "Failed to import any courses")
			task.CompleteTask(taskID, response)
		} else if len(errs) > 0 {
			task.SetTaskMessage(taskID, "Imported "+strconv.Itoa(len(importedCourses))+" courses with "+strconv.Itoa(len(errs))+" errors")
			task.CompleteTask(taskID, response)
		} else {
			task.SetTaskMessage(taskID, "Successfully imported "+strconv.Itoa(len(importedCourses))+" courses")
			task.CompleteTask(taskID, response)
		}
	}()

	// return task ID so client can check progress
	response := BatchImportStartResponse{
		Message: "Import started",
		TaskID:  taskID,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

// GetCourseProgress handles GET /api/courses/{id}/progress?user_id={uuid} - shows course progress for user
func (h *CourseHandler) GetCourseProgress(w http.ResponseWriter, r *http.Request) {
	// extract course ID from URL path
	pathParts := strings.Split(r.URL.Path, "/")
	if len(pathParts) < 4 {
		http.Error(w, "Invalid URL path", http.StatusBadRequest)
		return
	}

	courseIDStr := pathParts[3]
	courseID, err := uuid.Parse(courseIDStr)
	if err != nil {
		http.Error(w, "Invalid course ID format", http.StatusBadRequest)
		return
	}

	// get user ID from query params
	userIDStr := r.URL.Query().Get("user_id")
	if userIDStr == "" {
		http.Error(w, "user_id query parameter is required", http.StatusBadRequest)
		return
	}

	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		http.Error(w, "Invalid user ID format", http.StatusBadRequest)
		return
	}

	// calculate course progress
	progress, err := h.Service.CalculateCourseProgress(r.Context(), userID, courseID)
	if err != nil {
		log.Printf("Error calculating course progress: %v", err)
		http.Error(w, "Failed to calculate progress", http.StatusInternalServerError)
		return
	}

	response := CourseProgressResponse{
		Message: "Course progress calculated",
		Data:    progress,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

// GetModuleProgress handles GET /api/modules/{id}/progress?user_id={uuid} - shows module progress for user
func (h *CourseHandler) GetModuleProgress(w http.ResponseWriter, r *http.Request) {
	// extract module ID from URL path
	pathParts := strings.Split(r.URL.Path, "/")
	if len(pathParts) < 4 {
		http.Error(w, "Invalid URL path", http.StatusBadRequest)
		return
	}

	moduleIDStr := pathParts[3]
	moduleID, err := uuid.Parse(moduleIDStr)
	if err != nil {
		http.Error(w, "Invalid module ID format", http.StatusBadRequest)
		return
	}

	// get user ID from query params
	userIDStr := r.URL.Query().Get("user_id")
	if userIDStr == "" {
		http.Error(w, "user_id query parameter is required", http.StatusBadRequest)
		return
	}

	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		http.Error(w, "Invalid user ID format", http.StatusBadRequest)
		return
	}

	// calculate module progress
	progress, err := h.Service.CalculateModuleProgress(r.Context(), userID, moduleID)
	if err != nil {
		log.Printf("Error calculating module progress: %v", err)
		http.Error(w, "Failed to calculate progress", http.StatusInternalServerError)
		return
	}

	response := ModuleProgressResponse{
		Message: "Module progress calculated",
		Data:    progress,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

// UpdateContentProgress handles POST /api/content/{id}/progress - updates progress for content item
func (h *CourseHandler) UpdateContentProgress(w http.ResponseWriter, r *http.Request) {
	// extract content item ID from URL path
	pathParts := strings.Split(r.URL.Path, "/")
	if len(pathParts) < 4 {
		http.Error(w, "Invalid URL path", http.StatusBadRequest)
		return
	}

	contentIDStr := pathParts[3]
	contentID, err := uuid.Parse(contentIDStr)
	if err != nil {
		http.Error(w, "Invalid content ID format", http.StatusBadRequest)
		return
	}

	// parse request body
	type progressUpdate struct {
		UserID       uuid.UUID `json:"user_id"`
		ProgressPct  float32   `json:"progress_pct"`
		LastPosition int       `json:"last_position,omitempty"`
		Completed    bool      `json:"completed,omitempty"`
	}

	var update progressUpdate
	if err := json.NewDecoder(r.Body).Decode(&update); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// update progress
	err = h.Service.UpdateContentItemProgress(r.Context(), update.UserID, contentID, update.ProgressPct, update.LastPosition)
	if err != nil {
		log.Printf("Error updating content progress: %v", err)
		http.Error(w, "Failed to update progress", http.StatusInternalServerError)
		return
	}

	response := ProgressUpdateResponse{
		Message: "Progress updated successfully",
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

// MarkContentCompleted handles POST /api/content/{id}/complete - marks content as completed
func (h *CourseHandler) MarkContentCompleted(w http.ResponseWriter, r *http.Request) {
	// extract content item ID from URL path
	pathParts := strings.Split(r.URL.Path, "/")
	if len(pathParts) < 4 {
		http.Error(w, "Invalid URL path", http.StatusBadRequest)
		return
	}

	contentIDStr := pathParts[3]
	contentID, err := uuid.Parse(contentIDStr)
	if err != nil {
		http.Error(w, "Invalid content ID format", http.StatusBadRequest)
		return
	}

	// parse request body
	type completeRequest struct {
		UserID uuid.UUID `json:"user_id"`
	}

	var req completeRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// mark as completed
	err = h.Service.MarkContentItemCompleted(r.Context(), req.UserID, contentID)
	if err != nil {
		log.Printf("Error marking content completed: %v", err)
		http.Error(w, "Failed to mark as completed", http.StatusInternalServerError)
		return
	}

	response := ContentCompleteResponse{
		Message: "Content marked as completed",
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

// GetUserProgressSummary handles GET /api/users/{id}/progress - shows overall progress summary
func (h *CourseHandler) GetUserProgressSummary(w http.ResponseWriter, r *http.Request) {
	// extract user ID from URL path
	pathParts := strings.Split(r.URL.Path, "/")
	if len(pathParts) < 4 {
		http.Error(w, "Invalid URL path", http.StatusBadRequest)
		return
	}

	userIDStr := pathParts[3]
	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		http.Error(w, "Invalid user ID format", http.StatusBadRequest)
		return
	}

	// get progress summary
	summary, err := h.Service.GetUserProgressSummary(r.Context(), userID)
	if err != nil {
		log.Printf("Error getting progress summary: %v", err)
		http.Error(w, "Failed to get progress summary", http.StatusInternalServerError)
		return
	}

	response := UserProgressSummaryResponse{
		Message: "Progress summary retrieved",
		Data:    summary,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

// TODO: still need handlers for GetCourse, UpdateCourse, DeleteCourse
