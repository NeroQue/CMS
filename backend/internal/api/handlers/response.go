package handlers

import (
	"encoding/json"
	"log"
	"net/http"
)

// Common response structures for consistency across all handlers
type ErrorResponse struct {
	Message string `json:"message"`
	Success bool   `json:"success"`
}

type SuccessResponse struct {
	Message string      `json:"message"`
	Success bool        `json:"success"`
	Data    interface{} `json:"data,omitempty"`
}

// Helper functions for consistent response handling

// SendErrorResponse sends a consistent error response with logging
func SendErrorResponse(w http.ResponseWriter, message string, statusCode int, logMessage string, err error) {
	// Log the detailed error
	if err != nil {
		log.Printf("%s: %v", logMessage, err)
	} else {
		log.Printf("%s", logMessage)
	}

	// Set headers and status code
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)

	// Send structured error response
	response := ErrorResponse{
		Message: message,
		Success: false,
	}

	if encodeErr := json.NewEncoder(w).Encode(response); encodeErr != nil {
		log.Printf("Failed to encode error response: %v", encodeErr)
	}
}

// SendSuccessResponse sends a consistent success response with logging
func SendSuccessResponse(w http.ResponseWriter, message string, data interface{}, logMessage string) {
	// Log the success
	log.Printf("%s", logMessage)

	// Set headers and status code
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	// Send structured success response
	response := SuccessResponse{
		Message: message,
		Success: true,
		Data:    data,
	}

	if err := json.NewEncoder(w).Encode(response); err != nil {
		log.Printf("Failed to encode success response: %v", err)
		SendErrorResponse(w, "Failed to encode response", http.StatusInternalServerError, "JSON encoding error", err)
	}
}

// SendCreatedResponse sends a consistent response for created resources
func SendCreatedResponse(w http.ResponseWriter, message string, data interface{}, logMessage string) {
	// Log the success
	log.Printf("%s", logMessage)

	// Set headers and status code
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)

	// Send structured success response
	response := SuccessResponse{
		Message: message,
		Success: true,
		Data:    data,
	}

	if err := json.NewEncoder(w).Encode(response); err != nil {
		log.Printf("Failed to encode created response: %v", err)
		SendErrorResponse(w, "Failed to encode response", http.StatusInternalServerError, "JSON encoding error", err)
	}
}

// ValidateJSONBody validates and decodes JSON request body
func ValidateJSONBody(r *http.Request, dest interface{}) error {
	if r.Body == nil {
		return &ValidationError{Message: "Request body is required"}
	}

	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields() // Strict validation

	if err := decoder.Decode(dest); err != nil {
		return &ValidationError{Message: "Invalid JSON format: " + err.Error()}
	}

	return nil
}

// ValidationError represents validation errors
type ValidationError struct {
	Message string
}

func (e *ValidationError) Error() string {
	return e.Message
}
