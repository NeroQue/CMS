package handlers

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/NeroQue/course-management-backend/internal/services"
)

// Response structs for admin endpoints
type FactoryResetResponse struct {
	Message string `json:"message"`
	Success bool   `json:"success"`
}

type DatabaseStatsResponse struct {
	Message string         `json:"message"`
	Stats   map[string]int `json:"stats"`
}

// AdminHandler handles administrative operations
type AdminHandler struct {
	Service *services.AdminService // admin operations go through here
}

// NewAdminHandler creates handler with injected admin service
func NewAdminHandler(service *services.AdminService) *AdminHandler {
	return &AdminHandler{Service: service}
}

// FactoryReset handles POST /api/admin/factory-reset - clears all database data
func (h *AdminHandler) FactoryReset(w http.ResponseWriter, r *http.Request) {
	// basic validation - make sure this is a POST request
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	log.Println("Factory reset requested by user")

	// clear all data using the admin service
	err := h.Service.FactoryResetDatabase(r.Context())
	if err != nil {
		log.Printf("Error during factory reset: %v", err)
		response := FactoryResetResponse{
			Message: "Factory reset failed: " + err.Error(),
			Success: false,
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(response)
		return
	}

	response := FactoryResetResponse{
		Message: "Database factory reset completed successfully - all data cleared",
		Success: true,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}

// GetStats handles GET /api/admin/stats - shows basic database statistics
func (h *AdminHandler) GetStats(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// get database stats from admin service
	stats, err := h.Service.GetDatabaseStats(r.Context())
	if err != nil {
		log.Printf("Error getting database stats: %v", err)
		http.Error(w, "Failed to get database stats", http.StatusInternalServerError)
		return
	}

	response := DatabaseStatsResponse{
		Message: "Database statistics retrieved",
		Stats:   stats,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}
