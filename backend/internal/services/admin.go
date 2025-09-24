package services

import (
	"context"
	"fmt"
	"log"

	"github.com/NeroQue/course-management-backend/internal/database"
	"github.com/NeroQue/course-management-backend/pkg/session"
	"github.com/NeroQue/course-management-backend/pkg/task"
)

// AdminService handles administrative operations like factory reset
type AdminService struct {
	DB *database.Queries // database access
}

// NewAdminService creates admin service with database dependency
func NewAdminService(db *database.Queries) *AdminService {
	return &AdminService{
		DB: db,
	}
}

// FactoryResetDatabase clears all data from the database
func (s *AdminService) FactoryResetDatabase(ctx context.Context) error {
	log.Println("Starting factory reset - clearing all database data")

	// use the generated database method to clear all data
	err := s.DB.FactoryResetDatabase(ctx)
	if err != nil {
		return fmt.Errorf("failed to reset database: %w", err)
	}

	// clear any in-memory session data
	log.Println("Clearing session data")
	if err := session.ClearAllSessions(); err != nil {
		log.Printf("Warning: failed to clear sessions: %v", err)
		// don't fail the whole reset for this
	}

	// clear any running tasks since users will be logged out
	log.Println("Clearing task data")
	task.CleanupOldTasks(0) // clear all tasks regardless of age

	log.Println("Factory reset completed successfully")
	return nil
}

// GetDatabaseStats returns basic stats about database contents
func (s *AdminService) GetDatabaseStats(ctx context.Context) (map[string]int, error) {
	stats := make(map[string]int)

	// count profiles
	profiles, err := s.DB.GetAllProfiles(ctx)
	if err != nil {
		log.Printf("Warning: couldn't count profiles: %v", err)
		stats["profiles"] = -1
	} else {
		stats["profiles"] = len(profiles)
	}

	// count courses
	courses, err := s.DB.ListCourses(ctx)
	if err != nil {
		log.Printf("Warning: couldn't count courses: %v", err)
		stats["courses"] = -1
	} else {
		stats["courses"] = len(courses)
	}

	// TODO: could add counts for modules, content_items, sessions, progress records
	// but keeping it simple for now

	return stats, nil
}
