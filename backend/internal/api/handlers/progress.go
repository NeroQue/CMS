package handlers

import (
	"log"
	"net/http"
	"strings"

	"github.com/NeroQue/course-management-backend/internal/database"
	"github.com/NeroQue/course-management-backend/internal/services"
	"github.com/google/uuid"
)

// ProgressHandler processes progress-related HTTP requests
type ProgressHandler struct {
	Service *services.CourseService // handles all course business logic including progress
}

// NewProgressHandler creates handler with injected service
func NewProgressHandler(service *services.CourseService) *ProgressHandler {
	return &ProgressHandler{Service: service}
}

// GetCourseProgress handles GET /api/courses/{id}/progress?user_id={uuid} - shows course progress for user
// @Summary Get course progress
// @Description Get progress information for a specific course and user
// @Tags progress
// @Accept json
// @Produce json
// @Param id path string true "Course ID"
// @Param user_id query string true "User ID"
// @Success 200 {object} map[string]interface{} "Course progress information"
// @Failure 400 {object} map[string]interface{} "Invalid course or user ID"
// @Failure 500 {object} map[string]interface{} "Internal server error"
// @Router /courses/{id}/progress [get]
func (h *ProgressHandler) GetCourseProgress(w http.ResponseWriter, r *http.Request) {
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
// @Summary Get module progress
// @Description Get progress information for a specific module and user
// @Tags progress
// @Accept json
// @Produce json
// @Param id path string true "Module ID"
// @Param user_id query string true "User ID"
// @Success 200 {object} map[string]interface{} "Module progress information"
// @Failure 400 {object} map[string]interface{} "Invalid module or user ID"
// @Failure 500 {object} map[string]interface{} "Internal server error"
// @Router /modules/{id}/progress [get]
func (h *ProgressHandler) GetModuleProgress(w http.ResponseWriter, r *http.Request) {
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

// GetContentProgress handles GET /api/content/{id}/progress?user_id={uuid} - gets progress for content item
// @Summary Get content progress
// @Description Get progress information for a specific content item and user
// @Tags progress
// @Accept json
// @Produce json
// @Param id path string true "Content Item ID"
// @Param user_id query string true "User ID"
// @Success 200 {object} map[string]interface{} "Content progress information"
// @Failure 400 {object} map[string]interface{} "Invalid content or user ID"
// @Failure 404 {object} map[string]interface{} "Progress not found"
// @Failure 500 {object} map[string]interface{} "Internal server error"
// @Router /content/{id}/progress [get]
func (h *ProgressHandler) GetContentProgress(w http.ResponseWriter, r *http.Request) {
	log.Printf("Content progress requested from IP: %s", r.RemoteAddr)

	// extract content item ID from URL path
	pathParts := strings.Split(r.URL.Path, "/")
	if len(pathParts) < 4 {
		SendErrorResponse(w, "Invalid URL path format", http.StatusBadRequest,
			"Invalid URL path in content progress request", nil)
		return
	}

	contentIDStr := pathParts[3]
	contentID, err := uuid.Parse(contentIDStr)
	if err != nil {
		SendErrorResponse(w, "Invalid content ID format", http.StatusBadRequest,
			"Invalid content UUID in progress request", err)
		return
	}

	// get user ID from query params
	userIDStr := r.URL.Query().Get("user_id")
	if userIDStr == "" {
		SendErrorResponse(w, "user_id query parameter is required", http.StatusBadRequest,
			"Missing user_id parameter in content progress request", nil)
		return
	}

	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		SendErrorResponse(w, "Invalid user ID format", http.StatusBadRequest,
			"Invalid user UUID in progress request", err)
		return
	}

	log.Printf("Getting content progress for content %s and user %s", contentID.String(), userID.String())

	// get content progress from database
	progress, err := h.Service.DB.GetUserProgressByContentItem(r.Context(), database.GetUserProgressByContentItemParams{
		UserID:        userID,
		ContentItemID: contentID,
	})

	if err != nil {
		// If no progress found, return empty progress (not an error)
		if err.Error() == "sql: no rows in result set" {
			emptyProgress := map[string]interface{}{
				"user_id":         userID,
				"content_item_id": contentID,
				"completed":       false,
				"progress_pct":    0,
				"last_position":   0,
			}
			SendSuccessResponse(w, "No progress found", emptyProgress,
				"No progress record exists for this content item")
			return
		}

		SendErrorResponse(w, "Failed to get content progress", http.StatusInternalServerError,
			"Error retrieving content progress", err)
		return
	}

	SendSuccessResponse(w, "Content progress retrieved", progress,
		"Content progress retrieved and returned")
}

// UpdateContentProgress handles POST /api/content/{id}/progress - updates progress for content item
// @Summary Update content progress
// @Description Update progress information for a specific content item
// @Tags progress
// @Accept json
// @Produce json
// @Param id path string true "Content Item ID"
// @Param progress body object{user_id=string,progress_pct=number,last_position=number,completed=boolean} true "Progress update data"
// @Success 200 {object} map[string]interface{} "Progress updated successfully"
// @Failure 400 {object} map[string]interface{} "Invalid request format or missing fields"
// @Failure 500 {object} map[string]interface{} "Internal server error"
// @Router /content/{id}/progress [post]
func (h *ProgressHandler) UpdateContentProgress(w http.ResponseWriter, r *http.Request) {
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

	// fetch updated progress to return to client
	progress, err := h.Service.DB.GetUserProgressByContentItem(r.Context(), database.GetUserProgressByContentItemParams{
		UserID:        update.UserID,
		ContentItemID: contentID,
	})
	if err != nil {
		// Still send success even if we can't fetch progress
		SendSuccessResponse(w, "Progress updated successfully", nil,
			"Content progress updated successfully")
		return
	}

	SendSuccessResponse(w, "Progress updated successfully", progress,
		"Content progress updated successfully")
}

// MarkContentCompleted handles POST /api/content/{id}/complete - marks content as completed
// @Summary Mark content as completed
// @Description Mark a specific content item as completed for a user
// @Tags progress
// @Accept json
// @Produce json
// @Param id path string true "Content Item ID"
// @Param request body object{user_id=string} true "User ID"
// @Success 200 {object} map[string]interface{} "Content marked as completed"
// @Failure 400 {object} map[string]interface{} "Invalid request format or missing user ID"
// @Failure 500 {object} map[string]interface{} "Internal server error"
// @Router /content/{id}/complete [post]
func (h *ProgressHandler) MarkContentCompleted(w http.ResponseWriter, r *http.Request) {
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

	// fetch updated progress to return to client
	progress, err := h.Service.DB.GetUserProgressByContentItem(r.Context(), database.GetUserProgressByContentItemParams{
		UserID:        req.UserID,
		ContentItemID: contentID,
	})
	if err != nil {
		// Still send success even if we can't fetch progress
		SendSuccessResponse(w, "Content marked as completed", nil,
			"Content successfully marked as completed")
		return
	}

	SendSuccessResponse(w, "Content marked as completed", progress,
		"Content successfully marked as completed")
}

// ResetCourseProgress handles DELETE /api/courses/{id}/progress?user_id={uuid} - resets all progress for a course
// @Summary Reset course progress
// @Description Delete all progress records for a user in a specific course
// @Tags progress
// @Accept json
// @Produce json
// @Param id path string true "Course ID"
// @Param user_id query string true "User ID"
// @Success 200 {object} map[string]interface{} "Progress reset successfully"
// @Failure 400 {object} map[string]interface{} "Invalid course or user ID"
// @Failure 500 {object} map[string]interface{} "Internal server error"
// @Router /courses/{id}/progress [delete]
func (h *ProgressHandler) ResetCourseProgress(w http.ResponseWriter, r *http.Request) {
	log.Printf("Course progress reset requested from IP: %s", r.RemoteAddr)

	// extract course ID from URL path
	pathParts := strings.Split(r.URL.Path, "/")
	if len(pathParts) < 4 {
		SendErrorResponse(w, "Invalid URL path format", http.StatusBadRequest,
			"Invalid URL path in progress reset request", nil)
		return
	}

	courseIDStr := pathParts[3]
	courseID, err := uuid.Parse(courseIDStr)
	if err != nil {
		SendErrorResponse(w, "Invalid course ID format", http.StatusBadRequest,
			"Invalid course UUID in progress reset", err)
		return
	}

	// get user ID from query params
	userIDStr := r.URL.Query().Get("user_id")
	if userIDStr == "" {
		SendErrorResponse(w, "user_id query parameter is required", http.StatusBadRequest,
			"Missing user_id parameter in progress reset", nil)
		return
	}

	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		SendErrorResponse(w, "Invalid user ID format", http.StatusBadRequest,
			"Invalid user UUID in progress reset", err)
		return
	}

	log.Printf("Resetting progress for course %s and user %s", courseID.String(), userID.String())

	// get all content items in the course
	course, err := h.Service.GetCourse(r.Context(), courseID)
	if err != nil {
		SendErrorResponse(w, "Course not found", http.StatusNotFound,
			"Error retrieving course", err)
		return
	}

	// delete progress for each content item
	deletedCount := 0
	for _, module := range course.Modules {
		for _, contentItem := range module.ContentItems {
			err := h.Service.DB.DeleteUserProgress(r.Context(), database.DeleteUserProgressParams{
				UserID:        userID,
				ContentItemID: contentItem.ID,
			})
			if err != nil {
				log.Printf("Warning: failed to delete progress for content %s: %v", contentItem.ID, err)
			} else {
				deletedCount++
			}
		}
	}

	log.Printf("Reset %d progress records for course %s", deletedCount, courseID.String())

	responseData := map[string]interface{}{
		"deleted_count": deletedCount,
		"course_id":     courseID,
		"user_id":       userID,
	}

	SendSuccessResponse(w, "Course progress reset successfully", responseData,
		"All progress records deleted for this course")
}

// GetUserProgressSummary handles GET /api/users/{id}/progress - shows overall progress summary
// @Summary Get user progress summary
// @Description Get overall progress summary for a user across all courses
// @Tags progress
// @Accept json
// @Produce json
// @Param id path string true "User ID"
// @Success 200 {object} map[string]interface{} "User progress summary"
// @Failure 400 {object} map[string]interface{} "Invalid user ID"
// @Failure 500 {object} map[string]interface{} "Internal server error"
// @Router /users/{id}/progress [get]
func (h *ProgressHandler) GetUserProgressSummary(w http.ResponseWriter, r *http.Request) {
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
