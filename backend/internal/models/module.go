package models

import (
	"database/sql"

	"github.com/google/uuid"
)

// Module represents a section within a course
type Module struct {
	ID           uuid.UUID      `json:"id"`                      // unique identifier
	CourseID     uuid.UUID      `json:"course_id,omitempty"`     // which course this belongs to
	Title        string         `json:"title"`                   // module name
	Description  string         `json:"description,omitempty"`   // what this module covers
	RelativePath string         `json:"relative_path"`           // path relative to courses dir
	Order        int            `json:"order,omitempty"`         // position in course
	ContentItems []*ContentItem `json:"content_items,omitempty"` // actual content

	// timestamps
	CreatedAt sql.NullTime `json:"created_at,omitempty"`
	UpdatedAt sql.NullTime `json:"updated_at,omitempty"`
}

// CreateModuleInput is what we expect when creating a new module
type CreateModuleInput struct {
	CourseID     uuid.UUID `json:"course_id"`
	Title        string    `json:"title"`
	Description  string    `json:"description,omitempty"`
	RelativePath string    `json:"relative_path"`
	Order        int       `json:"order,omitempty"`
}
