package handlers

import (
	"log"
	"net/http"
	"strings"

	"github.com/NeroQue/course-management-backend/internal/models"
	"github.com/NeroQue/course-management-backend/internal/services"
	"github.com/NeroQue/course-management-backend/pkg/session"
	"github.com/google/uuid"
)

// ProfileHandler processes profile-related HTTP requests
type ProfileHandler struct {
	Service *services.ProfileService // business logic goes through here
}

// NewProfileHandler creates handler with injected service
func NewProfileHandler(service *services.ProfileService) *ProfileHandler {
	return &ProfileHandler{Service: service}
}

// List handles GET /api/profiles - returns all user profiles
// @Summary List all profiles
// @Description Get a list of all user profiles in the system
// @Tags profiles
// @Accept json
// @Produce json
// @Success 200 {object} map[string]interface{} "Successfully retrieved profiles"
// @Failure 500 {object} map[string]interface{} "Internal server error"
// @Router /profiles [get]
func (h *ProfileHandler) List(w http.ResponseWriter, r *http.Request) {
	log.Printf("Profile list requested from IP: %s", r.RemoteAddr)

	// get profiles from service layer
	profiles, err := h.Service.GetAllProfiles(r.Context())
	if err != nil {
		SendErrorResponse(w, "Failed to retrieve profiles", http.StatusInternalServerError,
			"Error retrieving profiles from database", err)
		return
	}

	SendSuccessResponse(w, "Profiles retrieved successfully", profiles,
		"Successfully retrieved and returned profile list")
}

// Get handles GET /api/profiles/{id} - returns a single profile
// @Summary Get a profile by ID
// @Description Get a single user profile with up-to-date gamification stats
// @Tags profiles
// @Accept json
// @Produce json
// @Param id path string true "Profile ID"
// @Success 200 {object} map[string]interface{} "Successfully retrieved profile"
// @Failure 400 {object} map[string]interface{} "Invalid profile ID format"
// @Failure 404 {object} map[string]interface{} "Profile not found"
// @Failure 500 {object} map[string]interface{} "Internal server error"
// @Router /profiles/{id} [get]
func (h *ProfileHandler) Get(w http.ResponseWriter, r *http.Request) {
	log.Printf("Profile get requested from IP: %s", r.RemoteAddr)

	// extract profile ID from URL path
	pathParts := strings.Split(r.URL.Path, "/")
	if len(pathParts) < 4 {
		SendErrorResponse(w, "Invalid URL path format", http.StatusBadRequest,
			"Invalid URL path in profile get request", nil)
		return
	}

	profileIDStr := pathParts[3]
	profileID, err := uuid.Parse(profileIDStr)
	if err != nil {
		SendErrorResponse(w, "Invalid profile ID format", http.StatusBadRequest,
			"Invalid UUID format in profile get request", err)
		return
	}

	log.Printf("Getting profile: %s", profileID.String())

	// get profile from service
	profile, err := h.Service.GetProfileByID(r.Context(), profileID)
	if err != nil {
		SendErrorResponse(w, "Profile not found", http.StatusNotFound,
			"Profile does not exist", err)
		return
	}

	SendSuccessResponse(w, "Profile retrieved successfully", profile,
		"Profile "+profileID.String()+" retrieved successfully")
}

// Create handles POST /api/profiles - makes new profile
// @Summary Create a new profile
// @Description Create a new user profile with the provided name
// @Tags profiles
// @Accept json
// @Produce json
// @Param profile body models.Profile true "Profile object with name field"
// @Success 201 {object} map[string]interface{} "Profile created successfully"
// @Failure 400 {object} map[string]interface{} "Invalid request format or missing required fields"
// @Failure 500 {object} map[string]interface{} "Internal server error"
// @Router /profiles [post]
func (h *ProfileHandler) Create(w http.ResponseWriter, r *http.Request) {
	log.Printf("Profile creation requested from IP: %s", r.RemoteAddr)

	// parse and validate the request body
	var profile models.Profile
	if err := ValidateJSONBody(r, &profile); err != nil {
		SendErrorResponse(w, "Invalid request format: "+err.Error(), http.StatusBadRequest,
			"Invalid JSON in profile creation request", err)
		return
	}

	// basic validation for required fields
	if strings.TrimSpace(profile.Name) == "" {
		SendErrorResponse(w, "Profile name is required", http.StatusBadRequest,
			"Profile creation attempted with empty name", nil)
		return
	}

	log.Printf("Creating new profile with name: %s", profile.Name)

	// use service to create profile
	createdProfile, err := h.Service.CreateProfile(r.Context(), profile)
	if err != nil {
		SendErrorResponse(w, "Failed to create profile", http.StatusInternalServerError,
			"Error creating profile in database", err)
		return
	}

	SendCreatedResponse(w, "Profile created successfully", createdProfile,
		"Profile created successfully with ID: "+createdProfile.ID.String())
}

// Update handles PUT /api/profiles - updates existing profile
// @Summary Update a profile
// @Description Update an existing profile's name
// @Tags profiles
// @Accept json
// @Produce json
// @Param updateRequest body object{user_id=string,new_name=string} true "Update request with user_id and new_name"
// @Success 200 {object} map[string]interface{} "Profile updated successfully"
// @Failure 400 {object} map[string]interface{} "Invalid request format or missing required fields"
// @Failure 500 {object} map[string]interface{} "Internal server error"
// @Router /profiles [put]
func (h *ProfileHandler) Update(w http.ResponseWriter, r *http.Request) {
	log.Printf("Profile update requested from IP: %s", r.RemoteAddr)

	// expect user ID and new name in request
	type updateRequest struct {
		UserID  uuid.UUID `json:"user_id"`
		NewName string    `json:"new_name"`
	}

	var req updateRequest
	if err := ValidateJSONBody(r, &req); err != nil {
		SendErrorResponse(w, "Invalid request format: "+err.Error(), http.StatusBadRequest,
			"Invalid JSON in profile update request", err)
		return
	}

	// validate required fields
	if req.UserID == uuid.Nil {
		SendErrorResponse(w, "User ID is required", http.StatusBadRequest,
			"Profile update attempted with missing user ID", nil)
		return
	}

	if strings.TrimSpace(req.NewName) == "" {
		SendErrorResponse(w, "New name is required and cannot be empty", http.StatusBadRequest,
			"Profile update attempted with empty name", nil)
		return
	}

	log.Printf("Updating profile %s with new name: %s", req.UserID.String(), req.NewName)

	// let service handle the update logic
	updatedProfile, err := h.Service.UpdateProfileName(r.Context(), req.UserID, req.NewName)
	if err != nil {
		SendErrorResponse(w, "Failed to update profile", http.StatusInternalServerError,
			"Error updating profile in database", err)
		return
	}

	SendSuccessResponse(w, "Profile updated successfully", updatedProfile,
		"Profile "+req.UserID.String()+" updated successfully")
}

// Delete handles DELETE /api/profiles - removes a profile
// @Summary Delete a profile
// @Description Delete a user profile by ID
// @Tags profiles
// @Accept json
// @Produce json
// @Param deleteRequest body object{user_id=string} true "Delete request with user_id"
// @Success 200 {object} map[string]interface{} "Profile deleted successfully"
// @Failure 400 {object} map[string]interface{} "Invalid request format or missing user ID"
// @Failure 500 {object} map[string]interface{} "Internal server error"
// @Router /profiles [delete]
func (h *ProfileHandler) Delete(w http.ResponseWriter, r *http.Request) {
	log.Printf("Profile deletion requested from IP: %s", r.RemoteAddr)

	type deleteRequest struct {
		UserID uuid.UUID `json:"user_id"`
	}

	var req deleteRequest
	if err := ValidateJSONBody(r, &req); err != nil {
		SendErrorResponse(w, "Invalid request format: "+err.Error(), http.StatusBadRequest,
			"Invalid JSON in profile deletion request", err)
		return
	}

	// validate required fields
	if req.UserID == uuid.Nil {
		SendErrorResponse(w, "User ID is required", http.StatusBadRequest,
			"Profile deletion attempted with missing user ID", nil)
		return
	}

	log.Printf("Deleting profile: %s", req.UserID.String())

	// service handles the actual deletion
	if err := h.Service.DeleteProfileByID(r.Context(), req.UserID); err != nil {
		SendErrorResponse(w, "Failed to delete profile", http.StatusInternalServerError,
			"Error deleting profile from database", err)
		return
	}

	SendSuccessResponse(w, "Profile deleted successfully", nil,
		"Profile "+req.UserID.String()+" deleted successfully")
}

// SelectProfile handles POST /api/profiles/{id}/select - sets active profile
// @Summary Select an active profile
// @Description Set a profile as the currently active user profile
// @Tags profiles
// @Accept json
// @Produce json
// @Param id path string true "Profile ID"
// @Success 200 {object} map[string]interface{} "Profile selected successfully"
// @Failure 400 {object} map[string]interface{} "Invalid profile ID format"
// @Failure 404 {object} map[string]interface{} "Profile not found"
// @Failure 500 {object} map[string]interface{} "Internal server error"
// @Router /profiles/{id}/select [post]
func (h *ProfileHandler) SelectProfile(w http.ResponseWriter, r *http.Request) {
	log.Printf("Profile selection requested from IP: %s", r.RemoteAddr)

	// extract profile ID from URL path like /api/profiles/{id}/select
	pathParts := strings.Split(r.URL.Path, "/")
	if len(pathParts) < 4 {
		SendErrorResponse(w, "Invalid URL path format", http.StatusBadRequest,
			"Invalid URL path in profile selection", nil)
		return
	}

	profileIDStr := pathParts[3]
	profileID, err := uuid.Parse(profileIDStr)
	if err != nil {
		SendErrorResponse(w, "Invalid profile ID format", http.StatusBadRequest,
			"Invalid UUID format in profile selection", err)
		return
	}

	log.Printf("Selecting profile: %s", profileID.String())

	// make sure profile actually exists
	_, err = h.Service.GetProfileByID(r.Context(), profileID)
	if err != nil {
		SendErrorResponse(w, "Profile not found", http.StatusNotFound,
			"Attempted to select non-existent profile", err)
		return
	}

	// set as current user in session
	session.SetCurrentUser(profileID)

	SendSuccessResponse(w, "Profile selected successfully", nil,
		"Profile "+profileID.String()+" selected as active")
}
