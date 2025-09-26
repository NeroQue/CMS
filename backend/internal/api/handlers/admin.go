package handlers

import (
	"log"
	"net/http"

	"github.com/NeroQue/course-management-backend/internal/services"
)

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
	log.Printf("Factory reset requested from IP: %s", r.RemoteAddr)

	// clear all data using the admin service
	err := h.Service.FactoryResetDatabase(r.Context())
	if err != nil {
		SendErrorResponse(w, "Factory reset failed", http.StatusInternalServerError,
			"Error during factory reset operation", err)
		return
	}

	SendSuccessResponse(w, "Database factory reset completed successfully - all data cleared",
		nil, "Factory reset completed successfully")
}

// GetStats handles GET /api/admin/stats - shows basic database statistics
func (h *AdminHandler) GetStats(w http.ResponseWriter, r *http.Request) {
	log.Printf("Database stats requested from IP: %s", r.RemoteAddr)

	// get database stats from admin service
	stats, err := h.Service.GetDatabaseStats(r.Context())
	if err != nil {
		SendErrorResponse(w, "Failed to get database statistics", http.StatusInternalServerError,
			"Error retrieving database statistics", err)
		return
	}

	SendSuccessResponse(w, "Database statistics retrieved successfully", stats,
		"Database statistics retrieved and returned to client")
}
