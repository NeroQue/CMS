package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/NeroQue/course-management-backend/internal/api"
	"github.com/NeroQue/course-management-backend/internal/database"
	"github.com/NeroQue/course-management-backend/pkg/parser"
	"github.com/NeroQue/course-management-backend/pkg/session"
	"github.com/NeroQue/course-management-backend/pkg/util"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"

	_ "github.com/NeroQue/course-management-backend/cmd/api/docs"
)

// @title Course Management System API
// @version 1.0
// @description API for managing courses, modules, profiles, and user progress in a learning management system
// @termsOfService http://swagger.io/terms/

// @contact.name API Support
// @contact.email support@example.com

// @license.name MIT
// @license.url https://opensource.org/licenses/MIT

// @host localhost:8080
// @BasePath /api
// @schemes http https

// @securityDefinitions.apikey SessionAuth
// @in cookie
// @name session_id

// main entry point - sets up everything and starts the server
func main() {
	// load .env file if it exists
	err := godotenv.Load(".env")
	if err != nil {
		log.Printf("Warning: Failed to load .env file: %s\n", err)
		// not a big deal - Docker will set these anyway
	}

	dbURL := os.Getenv("DB_URL")
	coursesDir := util.GetCoursesDirectory()

	// setup course parsing stuff
	courseParser := parser.NewCourseParser(coursesDir)
	if err := courseParser.ValidateBasePath(); err != nil {
		log.Printf("Warning: %v\n", err)
		log.Println("Course functionality may be limited")
	} else {
		log.Printf("Courses directory configured: %s\n", coursesDir)
	}

	// connect to postgres
	db, err := sql.Open("postgres", dbURL)
	if err != nil {
		log.Fatalf("Failed to connect to database: %s\n", err)
		return
	}
	defer db.Close()

	queries := database.New(db)
	session.Initialize(queries) // global session store - not ideal but works

	// wire everything together
	server := api.NewServer(db, courseParser)
	handler := server.EnableCORS(server) // needed for frontend requests

	fmt.Println("Starting server on :8080")
	// TODO: make port configurable via env var
	if err := http.ListenAndServe(":8080", handler); err != nil {
		log.Fatalf("Could not start server: %s\n", err)
	}
}
