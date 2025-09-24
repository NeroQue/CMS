package models

import (
	"database/sql"

	"github.com/google/uuid"
)

// Course represents a complete learning course
type Course struct {
	ID uuid.UUID `json:"id"` // unique identifier

	Title       string `json:"title"`                 // course name
	Description string `json:"description,omitempty"` // what the course is about

	Creator   string    `json:"creator,omitempty"`    // who added it
	CreatorID uuid.UUID `json:"creator_id,omitempty"` // creator's profile ID/the profile who added it

	// file path stuff - BasePath not stored in DB, just used during processing
	BasePath     string `json:"base_path,omitempty"`
	RelativePath string `json:"relative_path"` // path relative to courses dir

	Modules []*Module `json:"modules,omitempty"` // course content

	// timestamps
	CreatedAt sql.NullTime `json:"created_at,omitempty"`
	UpdatedAt sql.NullTime `json:"updated_at,omitempty"`
}

// CreateCourseInput is what we expect when creating a new course
type CreateCourseInput struct {
	Title        string    `json:"title"`
	Description  string    `json:"description,omitempty"`
	Creator      string    `json:"creator,omitempty"`
	CreatorID    uuid.UUID `json:"creator_id,omitempty"`
	BasePath     string    `json:"base_path,omitempty"`
	RelativePath string    `json:"relative_path"`
}

// CourseWithProgress shows course + how much user has completed
type CourseWithProgress struct {
	Course         *Course `json:"course"`
	CompletedItems int     `json:"completed_items"`
	TotalItems     int     `json:"total_items"`
	CompletionPct  float32 `json:"completion_pct"`
	LastAccessedAt *string `json:"last_accessed_at,omitempty"`
}

// TODO: add methods for validating course data, checking permissions, etc.
