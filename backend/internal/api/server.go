package api

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"time"

	"github.com/NeroQue/course-management-backend/internal/api/handlers"
	"github.com/NeroQue/course-management-backend/internal/database"
	"github.com/NeroQue/course-management-backend/internal/services"
	"github.com/NeroQue/course-management-backend/pkg/parser"
	"github.com/NeroQue/course-management-backend/pkg/task"
)

// Server holds all the app components together
type Server struct {
	DB *database.Queries // direct db access - probably should refactor this later

	Router *http.ServeMux // handles routing requests

	// handlers for different parts of the API
	ProfileHandler *handlers.ProfileHandler
	CourseHandler  *handlers.CourseHandler
	TaskHandler    *handlers.TaskHandler
	AdminHandler   *handlers.AdminHandler // for admin operations
}

// NewServer wires up all the dependencies and returns a ready-to-use server
func NewServer(db *sql.DB, courseParser *parser.CourseParser) *Server {
	dbQueries := database.New(db)

	task.Initialize()
	// start cleanup routine in background - cleans old tasks every hour
	go task.CleanupRoutine(1*time.Hour, 24*time.Hour)

	// create service layer instances
	profileSvc := services.NewProfileService(dbQueries)
	courseSvc := services.NewCourseService(dbQueries, courseParser)
	adminSvc := services.NewAdminService(dbQueries)

	// wire everything together
	server := &Server{
		DB:             dbQueries,
		Router:         http.NewServeMux(),
		ProfileHandler: handlers.NewProfileHandler(profileSvc),
		CourseHandler:  handlers.NewCourseHandler(courseSvc),
		TaskHandler:    handlers.NewTaskHandler(),
		AdminHandler:   handlers.NewAdminHandler(adminSvc),
	}

	server.setupRoutes()
	return server
}

// setupRoutes maps all the endpoints to handler functions
func (s *Server) setupRoutes() {
	s.Router.HandleFunc("/api", s.HelloHandler)

	// profile management
	s.Router.HandleFunc("GET /api/profiles", s.ProfileHandler.List)
	s.Router.HandleFunc("POST /api/profiles", s.ProfileHandler.Create)
	s.Router.HandleFunc("PUT /api/profiles", s.ProfileHandler.Update)
	s.Router.HandleFunc("DELETE /api/profiles", s.ProfileHandler.Delete)
	s.Router.HandleFunc("POST /api/profiles/{id}/select", s.ProfileHandler.SelectProfile)

	// course stuff
	s.Router.HandleFunc("GET /api/courses", s.CourseHandler.List)
	s.Router.HandleFunc("POST /api/courses", s.CourseHandler.Create)
	s.Router.HandleFunc("GET /api/courses/directories", s.CourseHandler.ListDirectories)
	s.Router.HandleFunc("GET /api/courses/scan", s.CourseHandler.ScanNewCourses)
	s.Router.HandleFunc("POST /api/courses/batch", s.CourseHandler.BatchImport)

	// progress tracking endpoints
	s.Router.HandleFunc("GET /api/courses/{id}/progress", s.CourseHandler.GetCourseProgress)
	s.Router.HandleFunc("GET /api/modules/{id}/progress", s.CourseHandler.GetModuleProgress)
	s.Router.HandleFunc("POST /api/content/{id}/progress", s.CourseHandler.UpdateContentProgress)
	s.Router.HandleFunc("POST /api/content/{id}/complete", s.CourseHandler.MarkContentCompleted)
	s.Router.HandleFunc("GET /api/users/{id}/progress", s.CourseHandler.GetUserProgressSummary)

	// admin endpoints
	s.Router.HandleFunc("POST /api/admin/factory-reset", s.AdminHandler.FactoryReset)
	s.Router.HandleFunc("GET /api/admin/stats", s.AdminHandler.GetStats)

	// task tracking
	s.Router.HandleFunc("GET /api/tasks", s.TaskHandler.GetTask)
	s.Router.HandleFunc("POST /api/tasks/cleanup", s.TaskHandler.CleanupTasks)
}

// ServeHTTP implements the http.Handler interface
// This allows the server to be used directly with http.ListenAndServe
func (s *Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	// Delegate to the router
	s.Router.ServeHTTP(w, r)
}

// HelloHandler is a simple handler for the base API endpoint
// This is kept at the server level as it doesn't require business logic
func (s *Server) HelloHandler(w http.ResponseWriter, r *http.Request) {
	type responseData struct {
		Message string `json:"message"`
	}

	response := responseData{Message: "Hello World from Go Backend with CORS!"}
	jsonResponse, _ := json.Marshal(response)
	w.Header().Set("Content-Type", "application/json")
	w.Write(jsonResponse)
}
