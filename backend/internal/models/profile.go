package models

import (
	"database/sql"
	"fmt"
	"time"

	"github.com/google/uuid"
)

// Profile represents a user in the system
type Profile struct {
	ID uuid.UUID `json:"id"` // unique identifier

	Name string `json:"name"` // display name

	// gamification stuff
	Experience int `json:"experience"` // XP points
	Gems       int `json:"gems"`       // special currency
	Streak     int `json:"streak"`     // consecutive active days

	LastActiveDate sql.NullTime `json:"last_active_date,omitempty"` // for streak tracking

	// timestamps
	CreatedAt sql.NullTime `json:"created_at,omitempty"`
	UpdatedAt sql.NullTime `json:"updated_at,omitempty"`
}

// CreateProfileInput is what we expect when creating a new profile
type CreateProfileInput struct {
	Name string `json:"name"`
}

// UpdateProfileInput is what we expect when updating a profile
type UpdateProfileInput struct {
	Name string `json:"name,omitempty"`
}

// GamificationUpdate represents changes to user's game stats
type GamificationUpdate struct {
	Experience int       `json:"experience"`
	Gems       int       `json:"gems"`
	Streak     int       `json:"streak"`
	LastActive time.Time `json:"last_active,omitempty"`
}

// String provides a string representation of the profile
// This is useful for logging and debugging
func (p *Profile) String() string {
	return fmt.Sprintf("Profile(ID=%s, Name=%s)", p.ID, p.Name)
}
