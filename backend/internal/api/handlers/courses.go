package handlers

import (
	"context"
	"log"
	"net/http"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/NeroQue/course-management-backend/internal/models"
	"github.com/NeroQue/course-management-backend/internal/services"
	"github.com/NeroQue/course-management-backend/pkg/session"
	"github.com/NeroQue/course-management-backend/pkg/task"
	"github.com/NeroQue/course-management-backend/pkg/util"
	"github.com/google/uuid"
)

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
	log.Printf("Course list requested from IP: %s", r.RemoteAddr)

	// get courses from service layer
	courses, err := h.Service.ListCourses(r.Context())
	if err != nil {
		SendErrorResponse(w, "Failed to retrieve courses", http.StatusInternalServerError,
			"Error retrieving courses from database", err)
		return
	}

	SendSuccessResponse(w, "Courses retrieved successfully", courses,
		"Successfully retrieved and returned course list")
}

// Create handles POST /api/courses - makes new course from directory
func (h *CourseHandler) Create(w http.ResponseWriter, r *http.Request) {
	log.Printf("Course creation requested from IP: %s", r.RemoteAddr)

	var input models.CreateCourseInput
	if err := ValidateJSONBody(r, &input); err != nil {
		SendErrorResponse(w, "Invalid request format: "+err.Error(), http.StatusBadRequest,
			"Invalid JSON in course creation request", err)
		return
	}

	// basic validation
	if strings.TrimSpace(input.Title) == "" {
		SendErrorResponse(w, "Course title is required", http.StatusBadRequest,
			"Course creation attempted with empty title", nil)
		return
	}

	if strings.TrimSpace(input.RelativePath) == "" {
		SendErrorResponse(w, "Relative path is required", http.StatusBadRequest,
			"Course creation attempted with empty relative path", nil)
		return
	}

	// need user logged in to create courses
	userID := session.GetCurrentUser()
	if userID == uuid.Nil {
		SendErrorResponse(w, "You must be logged in to create courses", http.StatusUnauthorized,
			"Unauthorized course creation attempt", nil)
		return
	}

	if input.BasePath == "" {
		input.BasePath = util.GetCoursesDirectory()
	}

	directoryPath := filepath.Join(input.BasePath, input.RelativePath)
	log.Printf("Creating course from directory: %s for user: %s", directoryPath, userID.String())

	// let service handle the actual import
	course, err := h.Service.ImportCourse(r.Context(), directoryPath, userID)
	if err != nil {
		SendErrorResponse(w, "Failed to create course: "+err.Error(), http.StatusBadRequest,
			"Error importing course from directory", err)
		return
	}

	SendCreatedResponse(w, "Course created successfully", course,
		"Course created successfully with ID: "+course.ID.String())
}

// ListDirectories handles GET /api/courses/directories - shows available dirs
func (h *CourseHandler) ListDirectories(w http.ResponseWriter, r *http.Request) {
	log.Printf("Course directories list requested from IP: %s", r.RemoteAddr)

	directories, err := h.Service.Parser.ListCourseDirectories()
	if err != nil {
		SendErrorResponse(w, "Failed to list directories", http.StatusInternalServerError,
			"Error listing course directories", err)
		return
	}

	SendSuccessResponse(w, "Directories retrieved successfully", directories,
		"Successfully retrieved course directories list")
}

// ScanNewCourses handles GET /api/courses/scan - finds dirs not imported yet
func (h *CourseHandler) ScanNewCourses(w http.ResponseWriter, r *http.Request) {
	log.Printf("New courses scan requested from IP: %s", r.RemoteAddr)

	// compare filesystem with database to find new ones
	newDirectories, err := h.Service.ScanNewCourses(r.Context())
	if err != nil {
		SendErrorResponse(w, "Failed to scan for new courses", http.StatusInternalServerError,
			"Error scanning for new courses", err)
		return
	}

	// Create custom response with count
	responseData := map[string]interface{}{
		"count":       len(newDirectories),
		"directories": newDirectories,
	}

	SendSuccessResponse(w, "New course directories found", responseData,
		"Found "+strconv.Itoa(len(newDirectories))+" new course directories")
}

// BatchImport handles POST /api/courses/batch - imports multiple courses at once
func (h *CourseHandler) BatchImport(w http.ResponseWriter, r *http.Request) {
	log.Printf("Batch course import requested from IP: %s", r.RemoteAddr)

	var request BatchImportRequest
	if err := ValidateJSONBody(r, &request); err != nil {
		SendErrorResponse(w, "Invalid request format: "+err.Error(), http.StatusBadRequest,
			"Invalid JSON in batch import request", err)
		return
	}

	if len(request.Courses) == 0 {
		SendErrorResponse(w, "No courses provided for import", http.StatusBadRequest,
			"Batch import attempted with empty course list", nil)
		return
	}

	userID := session.GetCurrentUser()
	if userID == uuid.Nil {
		SendErrorResponse(w, "You must be logged in to import courses", http.StatusUnauthorized,
			"Unauthorized batch import attempt", nil)
		return
	}

	// create background task since this might take a while
	taskID := task.CreateTask("batch_import")
	log.Printf("Starting batch import task %s for %d courses", taskID, len(request.Courses))

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
			log.Printf("Batch import %s failed completely", taskID)
		} else if len(errs) > 0 {
			task.SetTaskMessage(taskID, "Imported "+strconv.Itoa(len(importedCourses))+" courses with "+strconv.Itoa(len(errs))+" errors")
			task.CompleteTask(taskID, response)
			log.Printf("Batch import %s completed with partial success", taskID)
		} else {
			task.SetTaskMessage(taskID, "Successfully imported "+strconv.Itoa(len(importedCourses))+" courses")
			task.CompleteTask(taskID, response)
			log.Printf("Batch import %s completed successfully", taskID)
		}
	}()

	// return task ID so client can check progress
	responseData := map[string]string{"task_id": taskID}
	SendSuccessResponse(w, "Import started", responseData,
		"Batch import task created with ID: "+taskID)
}

// GetCourseProgress handles GET /api/courses/{id}/progress?user_id={uuid} - shows course progress for user
func (h *CourseHandler) GetCourseProgress(w http.ResponseWriter, r *http.Request) {
	log.Printf("Course progress requested from IP: %s", r.RemoteAddr)

	// extract course ID from URL path
	pathParts := strings.Split(r.URL.Path, "/")
	if len(pathParts) < 4 {
		SendErrorResponse(w, "Invalid URL path format", http.StatusBadRequest,
			"Invalid URL path in course progress request", nil)
		return
	}

	courseIDStr := pathParts[3]
	courseID, err := uuid.Parse(courseIDStr)
	if err != nil {
		SendErrorResponse(w, "Invalid course ID format", http.StatusBadRequest,
			"Invalid course UUID in progress request", err)
		return
	}

	// get user ID from query params
	userIDStr := r.URL.Query().Get("user_id")
	if userIDStr == "" {
		SendErrorResponse(w, "user_id query parameter is required", http.StatusBadRequest,
			"Missing user_id parameter in progress request", nil)
		return
	}

	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		SendErrorResponse(w, "Invalid user ID format", http.StatusBadRequest,
			"Invalid user UUID in progress request", err)
		return
	}

	log.Printf("Calculating course progress for course %s and user %s", courseID.String(), userID.String())

	// calculate course progress
	progress, err := h.Service.CalculateCourseProgress(r.Context(), userID, courseID)
	if err != nil {
		SendErrorResponse(w, "Failed to calculate progress", http.StatusInternalServerError,
			"Error calculating course progress", err)
		return
	}

	SendSuccessResponse(w, "Course progress calculated", progress,
		"Course progress calculated and returned")
}

// GetModuleProgress handles GET /api/modules/{id}/progress?user_id={uuid} - shows module progress for user
func (h *CourseHandler) GetModuleProgress(w http.ResponseWriter, r *http.Request) {
	log.Printf("Module progress requested from IP: %s", r.RemoteAddr)

	// extract module ID from URL path
	pathParts := strings.Split(r.URL.Path, "/")
	if len(pathParts) < 4 {
		SendErrorResponse(w, "Invalid URL path format", http.StatusBadRequest,
			"Invalid URL path in module progress request", nil)
		return
	}

	moduleIDStr := pathParts[3]
	moduleID, err := uuid.Parse(moduleIDStr)
	if err != nil {
		SendErrorResponse(w, "Invalid module ID format", http.StatusBadRequest,
			"Invalid module UUID in progress request", err)
		return
	}

	// get user ID from query params
	userIDStr := r.URL.Query().Get("user_id")
	if userIDStr == "" {
		SendErrorResponse(w, "user_id query parameter is required", http.StatusBadRequest,
			"Missing user_id parameter in progress request", nil)
		return
	}

	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		SendErrorResponse(w, "Invalid user ID format", http.StatusBadRequest,
			"Invalid user UUID in progress request", err)
		return
	}

	log.Printf("Calculating module progress for module %s and user %s", moduleID.String(), userID.String())

	// calculate module progress
	progress, err := h.Service.CalculateModuleProgress(r.Context(), userID, moduleID)
	if err != nil {
		SendErrorResponse(w, "Failed to calculate progress", http.StatusInternalServerError,
			"Error calculating module progress", err)
		return
	}

	SendSuccessResponse(w, "Module progress calculated", progress,
		"Module progress calculated and returned")
}

// UpdateContentProgress handles POST /api/content/{id}/progress - updates progress for content item
func (h *CourseHandler) UpdateContentProgress(w http.ResponseWriter, r *http.Request) {
	log.Printf("Content progress update requested from IP: %s", r.RemoteAddr)

	// extract content item ID from URL path
	pathParts := strings.Split(r.URL.Path, "/")
	if len(pathParts) < 4 {
		SendErrorResponse(w, "Invalid URL path format", http.StatusBadRequest,
			"Invalid URL path in content progress update", nil)
		return
	}

	contentIDStr := pathParts[3]
	contentID, err := uuid.Parse(contentIDStr)
	if err != nil {
		SendErrorResponse(w, "Invalid content ID format", http.StatusBadRequest,
			"Invalid content UUID in progress update", err)
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
	if err := ValidateJSONBody(r, &update); err != nil {
		SendErrorResponse(w, "Invalid request format: "+err.Error(), http.StatusBadRequest,
			"Invalid JSON in progress update request", err)
		return
	}

	// validate required fields
	if update.UserID == uuid.Nil {
		SendErrorResponse(w, "User ID is required", http.StatusBadRequest,
			"Progress update attempted with missing user ID", nil)
		return
	}

	log.Printf("Updating content progress for content %s, user %s, progress %.1f%%",
		contentID.String(), update.UserID.String(), update.ProgressPct)

	// update progress
	err = h.Service.UpdateContentItemProgress(r.Context(), update.UserID, contentID, update.ProgressPct, update.LastPosition)
	if err != nil {
		SendErrorResponse(w, "Failed to update progress", http.StatusInternalServerError,
			"Error updating content progress", err)
		return
	}

	SendSuccessResponse(w, "Progress updated successfully", nil,
		"Content progress updated successfully")
}

// MarkContentCompleted handles POST /api/content/{id}/complete - marks content as completed
func (h *CourseHandler) MarkContentCompleted(w http.ResponseWriter, r *http.Request) {
	log.Printf("Content completion requested from IP: %s", r.RemoteAddr)

	// extract content item ID from URL path
	pathParts := strings.Split(r.URL.Path, "/")
	if len(pathParts) < 4 {
		SendErrorResponse(w, "Invalid URL path format", http.StatusBadRequest,
			"Invalid URL path in content completion", nil)
		return
	}

	contentIDStr := pathParts[3]
	contentID, err := uuid.Parse(contentIDStr)
	if err != nil {
		SendErrorResponse(w, "Invalid content ID format", http.StatusBadRequest,
			"Invalid content UUID in completion request", err)
		return
	}

	// parse request body
	type completeRequest struct {
		UserID uuid.UUID `json:"user_id"`
	}

	var req completeRequest
	if err := ValidateJSONBody(r, &req); err != nil {
		SendErrorResponse(w, "Invalid request format: "+err.Error(), http.StatusBadRequest,
			"Invalid JSON in completion request", err)
		return
	}

	// validate required fields
	if req.UserID == uuid.Nil {
		SendErrorResponse(w, "User ID is required", http.StatusBadRequest,
			"Content completion attempted with missing user ID", nil)
		return
	}

	log.Printf("Marking content %s as completed for user %s", contentID.String(), req.UserID.String())

	// mark as completed
	err = h.Service.MarkContentItemCompleted(r.Context(), req.UserID, contentID)
	if err != nil {
		SendErrorResponse(w, "Failed to mark as completed", http.StatusInternalServerError,
			"Error marking content as completed", err)
		return
	}

	SendSuccessResponse(w, "Content marked as completed", nil,
		"Content successfully marked as completed")
}

// GetUserProgressSummary handles GET /api/users/{id}/progress - shows overall progress summary
func (h *CourseHandler) GetUserProgressSummary(w http.ResponseWriter, r *http.Request) {
	log.Printf("User progress summary requested from IP: %s", r.RemoteAddr)

	// extract user ID from URL path
	pathParts := strings.Split(r.URL.Path, "/")
	if len(pathParts) < 4 {
		SendErrorResponse(w, "Invalid URL path format", http.StatusBadRequest,
			"Invalid URL path in progress summary request", nil)
		return
	}

	userIDStr := pathParts[3]
	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		SendErrorResponse(w, "Invalid user ID format", http.StatusBadRequest,
			"Invalid user UUID in progress summary request", err)
		return
	}

	log.Printf("Getting progress summary for user %s", userID.String())

	// get progress summary
	summary, err := h.Service.GetUserProgressSummary(r.Context(), userID)
	if err != nil {
		SendErrorResponse(w, "Failed to get progress summary", http.StatusInternalServerError,
			"Error getting user progress summary", err)
		return
	}

	SendSuccessResponse(w, "Progress summary retrieved", summary,
		"User progress summary retrieved and returned")
}
