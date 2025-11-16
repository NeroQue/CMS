package handlers

import (
	"context"
	"log"
	"net/http"
	"os"
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
// @Summary List all courses
// @Description Get a list of all available courses
// @Tags courses
// @Accept json
// @Produce json
// @Success 200 {object} map[string]interface{} "Successfully retrieved courses"
// @Failure 500 {object} map[string]interface{} "Internal server error"
// @Router /courses [get]
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

// Get handles GET /api/courses/{id} - returns single course with modules and content
// @Summary Get a course by ID
// @Description Get detailed information about a specific course including modules and content items
// @Tags courses
// @Accept json
// @Produce json
// @Param id path string true "Course ID"
// @Success 200 {object} map[string]interface{} "Successfully retrieved course"
// @Failure 400 {object} map[string]interface{} "Invalid course ID"
// @Failure 404 {object} map[string]interface{} "Course not found"
// @Failure 500 {object} map[string]interface{} "Internal server error"
// @Router /courses/{id} [get]
func (h *CourseHandler) Get(w http.ResponseWriter, r *http.Request) {
	log.Printf("Course detail requested from IP: %s", r.RemoteAddr)

	// extract course ID from URL path
	pathParts := strings.Split(r.URL.Path, "/")
	if len(pathParts) < 4 {
		SendErrorResponse(w, "Invalid URL path format", http.StatusBadRequest,
			"Invalid URL path in course detail request", nil)
		return
	}

	courseIDStr := pathParts[3]
	courseID, err := uuid.Parse(courseIDStr)
	if err != nil {
		SendErrorResponse(w, "Invalid course ID format", http.StatusBadRequest,
			"Invalid course UUID in detail request", err)
		return
	}

	log.Printf("Getting course details for course %s", courseID.String())

	// get course with modules and content items
	course, err := h.Service.GetCourse(r.Context(), courseID)
	if err != nil {
		SendErrorResponse(w, "Failed to get course", http.StatusNotFound,
			"Error getting course details", err)
		return
	}

	SendSuccessResponse(w, "Course retrieved successfully", course,
		"Course details retrieved and returned")
}

// Create handles POST /api/courses - makes new course from directory
// @Summary Create a new course
// @Description Import a new course from a directory path
// @Tags courses
// @Accept json
// @Produce json
// @Param course body models.CreateCourseInput true "Course creation input with title and relative_path"
// @Success 201 {object} map[string]interface{} "Course created successfully"
// @Failure 400 {object} map[string]interface{} "Invalid request format or missing required fields"
// @Failure 401 {object} map[string]interface{} "User must be logged in"
// @Failure 500 {object} map[string]interface{} "Internal server error"
// @Security SessionAuth
// @Router /courses [post]
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
// @Summary List course directories
// @Description Get a list of all available course directories in the filesystem
// @Tags courses
// @Accept json
// @Produce json
// @Success 200 {object} map[string]interface{} "Successfully retrieved directories"
// @Failure 500 {object} map[string]interface{} "Internal server error"
// @Router /courses/directories [get]
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
// @Summary Scan for new courses
// @Description Find course directories that haven't been imported yet
// @Tags courses
// @Accept json
// @Produce json
// @Success 200 {object} map[string]interface{} "Successfully scanned for new courses"
// @Failure 500 {object} map[string]interface{} "Internal server error"
// @Router /courses/scan [get]
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
// @Summary Batch import courses
// @Description Import multiple courses at once from provided course inputs
// @Tags courses
// @Accept json
// @Produce json
// @Param batchRequest body BatchImportRequest true "Batch import request with array of courses"
// @Success 200 {object} map[string]interface{} "Returns task ID for tracking import progress"
// @Failure 400 {object} map[string]interface{} "Invalid request format or empty course list"
// @Failure 401 {object} map[string]interface{} "User must be logged in"
// @Failure 500 {object} map[string]interface{} "Internal server error"
// @Security SessionAuth
// @Router /courses/batch [post]
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

// ServeContentFile handles GET /api/content/{id}/file - serves the actual content file
// @Summary Serve content file
// @Description Serve the actual file for a content item (video, PDF, etc.)
// @Tags content
// @Produce octet-stream
// @Param id path string true "Content Item ID"
// @Success 200 {file} binary "Content file"
// @Failure 400 {object} map[string]interface{} "Invalid content ID"
// @Failure 404 {object} map[string]interface{} "File not found"
// @Failure 500 {object} map[string]interface{} "Internal server error"
// @Router /content/{id}/file [get]
func (h *CourseHandler) ServeContentFile(w http.ResponseWriter, r *http.Request) {
	log.Printf("Content file requested from IP: %s", r.RemoteAddr)

	// extract content item ID from URL path
	pathParts := strings.Split(r.URL.Path, "/")
	if len(pathParts) < 4 {
		SendErrorResponse(w, "Invalid URL path format", http.StatusBadRequest,
			"Invalid URL path in file serve request", nil)
		return
	}

	contentIDStr := pathParts[3]
	contentID, err := uuid.Parse(contentIDStr)
	if err != nil {
		SendErrorResponse(w, "Invalid content ID format", http.StatusBadRequest,
			"Invalid content UUID in file serve request", err)
		return
	}

	log.Printf("Serving file for content %s", contentID.String())

	// get content item from database to get file path
	contentItem, err := h.Service.DB.GetContentItem(r.Context(), contentID)
	if err != nil {
		SendErrorResponse(w, "Content item not found", http.StatusNotFound,
			"Error retrieving content item", err)
		return
	}

	// construct full file path
	fullPath := filepath.Join(h.Service.Parser.BasePath, contentItem.RelativePath)

	// Docker container path adjustment
	if strings.HasPrefix(fullPath, "/courses/") {
		adjustedPath := filepath.Join("../", fullPath)
		if _, err := os.Stat(adjustedPath); err == nil {
			fullPath = adjustedPath
		}
	}

	log.Printf("Serving file from path: %s", fullPath)

	// check if file exists
	if _, err := os.Stat(fullPath); err != nil {
		SendErrorResponse(w, "File not found", http.StatusNotFound,
			"Content file does not exist on disk", err)
		return
	}

	// serve the file
	http.ServeFile(w, r, fullPath)
}

// Delete handles DELETE /api/courses/{id} - removes a course from the database
// @Summary Delete a course
// @Description Delete a course from the database (does not remove files from disk)
// @Tags courses
// @Accept json
// @Produce json
// @Param id path string true "Course ID"
// @Success 200 {object} map[string]interface{} "Course deleted successfully"
// @Failure 400 {object} map[string]interface{} "Invalid course ID"
// @Failure 404 {object} map[string]interface{} "Course not found"
// @Failure 500 {object} map[string]interface{} "Internal server error"
// @Router /courses/{id} [delete]
func (h *CourseHandler) Delete(w http.ResponseWriter, r *http.Request) {
	log.Printf("Course deletion requested from IP: %s", r.RemoteAddr)

	// extract course ID from URL path
	pathParts := strings.Split(r.URL.Path, "/")
	if len(pathParts) < 4 {
		SendErrorResponse(w, "Invalid URL path format", http.StatusBadRequest,
			"Invalid URL path in course deletion request", nil)
		return
	}

	courseIDStr := pathParts[3]
	courseID, err := uuid.Parse(courseIDStr)
	if err != nil {
		SendErrorResponse(w, "Invalid course ID format", http.StatusBadRequest,
			"Invalid course UUID in deletion request", err)
		return
	}

	log.Printf("Deleting course %s", courseID.String())

	// delete the course using the service
	err = h.Service.DeleteCourse(r.Context(), courseID)
	if err != nil {
		SendErrorResponse(w, "Failed to delete course", http.StatusInternalServerError,
			"Error deleting course from database", err)
		return
	}

	SendSuccessResponse(w, "Course deleted successfully", nil,
		"Course deleted from database: "+courseID.String())
}

// CheckCourseExists handles GET /api/courses/{id}/exists - checks if course directory exists on disk
// @Summary Check if course exists on disk
// @Description Verify that the course directory and files still exist on the filesystem
// @Tags courses
// @Accept json
// @Produce json
// @Param id path string true "Course ID"
// @Success 200 {object} map[string]interface{} "Returns exists: true/false and missing files if any"
// @Failure 400 {object} map[string]interface{} "Invalid course ID"
// @Failure 404 {object} map[string]interface{} "Course not found in database"
// @Failure 500 {object} map[string]interface{} "Internal server error"
// @Router /courses/{id}/exists [get]
func (h *CourseHandler) CheckCourseExists(w http.ResponseWriter, r *http.Request) {
	log.Printf("Course existence check requested from IP: %s", r.RemoteAddr)

	// extract course ID from URL path
	pathParts := strings.Split(r.URL.Path, "/")
	if len(pathParts) < 4 {
		SendErrorResponse(w, "Invalid URL path format", http.StatusBadRequest,
			"Invalid URL path in existence check request", nil)
		return
	}

	courseIDStr := pathParts[3]
	courseID, err := uuid.Parse(courseIDStr)
	if err != nil {
		SendErrorResponse(w, "Invalid course ID format", http.StatusBadRequest,
			"Invalid course UUID in existence check request", err)
		return
	}

	log.Printf("Checking existence for course %s", courseID.String())

	// check if course exists on disk
	exists, missingPaths, err := h.Service.CheckCourseExistsOnDisk(r.Context(), courseID)
	if err != nil {
		SendErrorResponse(w, "Failed to check course existence", http.StatusInternalServerError,
			"Error checking course existence", err)
		return
	}

	responseData := map[string]interface{}{
		"exists":        exists,
		"missing_paths": missingPaths,
	}

	SendSuccessResponse(w, "Course existence checked", responseData,
		"Course existence verification completed")
}
