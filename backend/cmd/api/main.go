package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/NeroQue/course-management-backend/internal/api"
	"github.com/NeroQue/course-management-backend/internal/database"
	"github.com/NeroQue/course-management-backend/pkg/parser"
	"github.com/NeroQue/course-management-backend/pkg/session"
	"github.com/NeroQue/course-management-backend/pkg/util"
	"github.com/NeroQue/course-management-backend/sql/schema"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
	"github.com/pressly/goose/v3"

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

	// connect to postgres with retry - gives the DB container time to start
	db, err := connectWithRetry(dbURL, 30, 2*time.Second)
	if err != nil {
		log.Fatalf("Failed to connect to database after retries: %s\n", err)
		return
	}
	defer db.Close()

	// run migrations before anything else touches the DB
	if err := runMigrations(db); err != nil {
		log.Fatalf("Failed to run database migrations: %s\n", err)
		return
	}

	queries := database.New(db)
	session.Initialize(queries) // global session store - not ideal but works

	// wire everything together
	server := api.NewServer(db, courseParser)
	handler := server.EnableCORS(server) // needed for frontend requests

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	fmt.Printf("Starting server on :%s\n", port)
	if err := http.ListenAndServe(":"+port, handler); err != nil {
		log.Fatalf("Could not start server: %s\n", err)
	}
}

// connectWithRetry tries to connect to the database, retrying on failure.
// This handles the Docker race condition where the backend starts before Postgres is ready.
func connectWithRetry(dbURL string, maxRetries int, delay time.Duration) (*sql.DB, error) {
	var db *sql.DB
	var err error

	for i := 0; i < maxRetries; i++ {
		db, err = sql.Open("postgres", dbURL)
		if err != nil {
			log.Printf("Failed to open database (attempt %d/%d): %s", i+1, maxRetries, err)
			time.Sleep(delay)
			continue
		}

		err = db.Ping()
		if err == nil {
			log.Println("Successfully connected to database")
			return db, nil
		}

		log.Printf("Database not ready (attempt %d/%d): %s", i+1, maxRetries, err)
		db.Close()
		time.Sleep(delay)
	}

	return nil, fmt.Errorf("could not connect to database after %d attempts: %w", maxRetries, err)
}

// runMigrations applies all pending goose migrations using the embedded SQL files.
func runMigrations(db *sql.DB) error {
	goose.SetBaseFS(schema.Migrations)

	if err := goose.SetDialect("postgres"); err != nil {
		return fmt.Errorf("failed to set goose dialect: %w", err)
	}

	if err := goose.Up(db, "."); err != nil {
		return fmt.Errorf("failed to run migrations: %w", err)
	}

	log.Println("Database migrations applied successfully")
	return nil
}
