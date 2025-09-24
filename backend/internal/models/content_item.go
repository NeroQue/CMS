package models

import (
	"database/sql"

	"github.com/google/uuid"
)

// ContentItem represents individual content like videos, PDFs, etc.
type ContentItem struct {
	ID       uuid.UUID `json:"id"`                  // unique identifier
	ModuleID uuid.UUID `json:"module_id,omitempty"` // which module this belongs to

	Title       string `json:"title"`                 // content name
	Description string `json:"description,omitempty"` // what this content is about

	RelativePath string `json:"relative_path"` // path to the actual file
	ContentType  string `json:"content_type"`  // video, pdf, text, etc.

	Duration int   `json:"duration,omitempty"` // seconds (for videos)
	Size     int64 `json:"size,omitempty"`     // file size in bytes
	Order    int   `json:"order,omitempty"`    // position in module

	// timestamps
	CreatedAt sql.NullTime `json:"created_at,omitempty"`
	UpdatedAt sql.NullTime `json:"updated_at,omitempty"`
}

// CreateContentItemInput is what we expect when creating new content
type CreateContentItemInput struct {
	ModuleID     uuid.UUID `json:"module_id"`
	Title        string    `json:"title"`
	Description  string    `json:"description,omitempty"`
	RelativePath string    `json:"relative_path"`
	ContentType  string    `json:"content_type"`
	Duration     int       `json:"duration,omitempty"`
	Size         int64     `json:"size,omitempty"`
	Order        int       `json:"order,omitempty"`
}
