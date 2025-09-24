package models

import (
	"database/sql"
	"time"

	"github.com/google/uuid"
)

// UserProgress tracks how far a user has gotten through content
type UserProgress struct {
	ID uuid.UUID `json:"id"` // unique identifier

	UserID        uuid.UUID `json:"user_id"`         // which user
	ContentItemID uuid.UUID `json:"content_item_id"` // which content they're working on

	Completed   bool    `json:"completed"`    // whether they finished it
	ProgressPct float32 `json:"progress_pct"` // how much done (0-100)

	LastPosition int          `json:"last_position,omitempty"` // seconds (for videos)
	LastAccessed sql.NullTime `json:"last_accessed,omitempty"` // when they last viewed it

	// timestamps
	CreatedAt sql.NullTime `json:"created_at,omitempty"`
	UpdatedAt sql.NullTime `json:"updated_at,omitempty"`
}

// CreateProgressInput is what we expect when tracking progress
type CreateProgressInput struct {
	UserID        uuid.UUID `json:"user_id"`
	ContentItemID uuid.UUID `json:"content_item_id"`
	Completed     bool      `json:"completed"`
	ProgressPct   float32   `json:"progress_pct"`
	LastPosition  int       `json:"last_position,omitempty"`
}

// ModuleProgress represents calculated progress for a module
type ModuleProgress struct {
	ModuleID       uuid.UUID  `json:"module_id"`
	UserID         uuid.UUID  `json:"user_id"`
	CompletedItems int        `json:"completed_items"`
	TotalItems     int        `json:"total_items"`
	CompletionPct  float32    `json:"completion_pct"`
	LastAccessedAt *time.Time `json:"last_accessed_at,omitempty"`
	IsCompleted    bool       `json:"is_completed"` // true when all content items done
}

// CourseProgress represents calculated progress for an entire course
type CourseProgress struct {
	CourseID          uuid.UUID  `json:"course_id"`
	UserID            uuid.UUID  `json:"user_id"`
	CompletedModules  int        `json:"completed_modules"`
	TotalModules      int        `json:"total_modules"`
	CompletedItems    int        `json:"completed_items"`
	TotalItems        int        `json:"total_items"`
	CompletionPct     float32    `json:"completion_pct"`
	LastAccessedAt    *time.Time `json:"last_accessed_at,omitempty"`
	IsCompleted       bool       `json:"is_completed"`                  // true when all modules done
	EstimatedTimeLeft int        `json:"estimated_time_left,omitempty"` // minutes
}

// ProgressSummary gives overall user progress across all courses
type ProgressSummary struct {
	UserID            uuid.UUID `json:"user_id"`
	TotalCourses      int       `json:"total_courses"`
	CompletedCourses  int       `json:"completed_courses"`
	InProgressCourses int       `json:"in_progress_courses"`
	TotalTimeSpent    int       `json:"total_time_spent"` // minutes
	StreakDays        int       `json:"streak_days"`
}
