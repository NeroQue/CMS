package handlers

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"

	"github.com/NeroQue/course-management-backend/internal/models"
	"github.com/NeroQue/course-management-backend/internal/services"
	"github.com/NeroQue/course-management-backend/pkg/session"
	"github.com/google/uuid"
)

// Response structs for profile endpoints
type ProfileListResponse struct {
	Message string           `json:"message"`
	Data    []models.Profile `json:"data"`
}

type ProfileResponse struct {
	Message string         `json:"message"`
	Data    models.Profile `json:"data"`
}

type ProfileDeleteResponse struct {
	Message string `json:"message"`
}

type ProfileSelectResponse struct {
	Message string `json:"message"`
}

// ProfileHandler processes profile-related HTTP requests
type ProfileHandler struct {
	Service *services.ProfileService // business logic goes through here
}

// NewProfileHandler creates handler with injected service
func NewProfileHandler(service *services.ProfileService) *ProfileHandler {
	return &ProfileHandler{Service: service}
}

// List handles GET /api/profiles - returns all user profiles
func (h *ProfileHandler) List(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	w.Header().Set("Content-Type", "application/json")

	// get profiles from service layer
	profiles, err := h.Service.GetAllProfiles(r.Context())
	if err != nil {
		log.Printf("Error retrieving profiles: %v", err)
		http.Error(w, "Failed to retrieve profiles", http.StatusInternalServerError)
		return
	}

	response := ProfileListResponse{
		Message: "Profiles retrieved successfully",
		Data:    profiles,
	}

	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(response); err != nil {
		log.Printf("Error encoding JSON response: %v", err)
		http.Error(w, "Failed to encode response", http.StatusInternalServerError)
		return
	}
}

// Create handles POST /api/profiles - makes new profile
func (h *ProfileHandler) Create(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	w.Header().Set("Content-Type", "application/json")

	// parse the request body
	var profile models.Profile
	if err := json.NewDecoder(r.Body).Decode(&profile); err != nil {
		log.Printf("Error decoding request body: %v", err)
		http.Error(w, "Invalid request format", http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	// use service to create profile
	createdProfile, err := h.Service.CreateProfile(r.Context(), profile)
	if err != nil {
		log.Printf("Error creating profile: %v", err)
		http.Error(w, "Failed to create profile", http.StatusInternalServerError)
		return
	}

	response := ProfileResponse{
		Message: "Profile created successfully",
		Data:    createdProfile,
	}

	w.WriteHeader(http.StatusCreated)
	if err := json.NewEncoder(w).Encode(response); err != nil {
		log.Printf("Error encoding response: %v", err)
		http.Error(w, "Failed to encode response", http.StatusInternalServerError)
		return
	}
}

// Update handles PUT /api/profiles - updates existing profile
func (h *ProfileHandler) Update(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPut {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	w.Header().Set("Content-Type", "application/json")

	// expect current name and new name in request
	type updateRequest struct {
		CurrentName string `json:"current_name"`
		NewName     string `json:"new_name"`
	}

	var req updateRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		log.Printf("Error decoding request body: %v", err)
		http.Error(w, "Invalid request format", http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	// let service handle the update logic
	updatedProfile, err := h.Service.UpdateProfileName(r.Context(), req.CurrentName, req.NewName)
	if err != nil {
		log.Printf("Error updating profile: %v", err)
		http.Error(w, "Failed to update profile", http.StatusInternalServerError)
		return
	}

	response := ProfileResponse{
		Message: "Profile updated successfully",
		Data:    updatedProfile,
	}

	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(response); err != nil {
		log.Printf("Error encoding response: %v", err)
		http.Error(w, "Failed to encode response", http.StatusInternalServerError)
		return
	}
}

// Delete handles DELETE /api/profiles - removes a profile
func (h *ProfileHandler) Delete(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodDelete {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	w.Header().Set("Content-Type", "application/json")

	type deleteRequest struct {
		Name string `json:"name"`
	}

	var req deleteRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		log.Printf("Error decoding request body: %v", err)
		http.Error(w, "Invalid request format", http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	// service handles the actual deletion
	if err := h.Service.DeleteProfileByName(r.Context(), req.Name); err != nil {
		log.Printf("Error deleting profile: %v", err)
		http.Error(w, "Failed to delete profile", http.StatusInternalServerError)
		return
	}

	response := ProfileDeleteResponse{
		Message: fmt.Sprintf("Profile %s deleted successfully", req.Name),
	}

	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(response); err != nil {
		log.Printf("Error encoding response: %v", err)
		return
	}
}

// SelectProfile handles POST /api/profiles/{id}/select - sets active profile
func (h *ProfileHandler) SelectProfile(w http.ResponseWriter, r *http.Request) {
	// extract profile ID from URL path like /api/profiles/{id}/select
	pathParts := strings.Split(r.URL.Path, "/")
	if len(pathParts) < 4 {
		http.Error(w, "Invalid URL path", http.StatusBadRequest)
		return
	}

	profileIDStr := pathParts[3]
	profileID, err := uuid.Parse(profileIDStr)
	if err != nil {
		http.Error(w, "Invalid profile ID format", http.StatusBadRequest)
		return
	}

	// make sure profile actually exists
	_, err = h.Service.GetProfileByID(r.Context(), profileID)
	if err != nil {
		log.Printf("Error retrieving profile: %v", err)
		http.Error(w, "Profile not found", http.StatusNotFound)
		return
	}

	// set as current user in session
	session.SetCurrentUser(profileID)

	response := ProfileSelectResponse{
		Message: "Profile selected successfully",
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}
