package api

import (
	"net/http"
)

// EnableCORS adds CORS headers so frontend can talk to the API
func (s *Server) EnableCORS(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// allow all origins for now - should probably restrict this later
		w.Header().Set("Access-Control-Allow-Origin", "*")

		// allow the HTTP methods we use
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")

		// need this for JSON requests
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

		// handle preflight requests from browser
		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusOK)
			return
		}

		// pass request along to actual handler
		next.ServeHTTP(w, r)
	})
}

// TODO: need to add middleware for auth, logging, etc.
