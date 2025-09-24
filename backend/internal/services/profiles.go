package services

import (
	"context"
	"errors"
	"fmt"
	"log"
	"strings"

	"github.com/NeroQue/course-management-backend/internal/database"
	"github.com/NeroQue/course-management-backend/internal/models"
	"github.com/google/uuid"
)

// ProfileService handles all the profile business logic
type ProfileService struct {
	DB *database.Queries // database access layer
}

// NewProfileService creates service with db dependency
func NewProfileService(db *database.Queries) *ProfileService {
	return &ProfileService{
		DB: db,
	}
}

// GetAllProfiles fetches all profiles from database
func (s *ProfileService) GetAllProfiles(ctx context.Context) ([]models.Profile, error) {
	profiles, err := s.DB.GetAllProfiles(ctx)
	if err != nil {
		log.Printf("Error retrieving profiles: %v", err)
		return nil, fmt.Errorf("failed to retrieve profiles: %w", err)
	}

	// convert db models to app models
	modelProfiles := make([]models.Profile, len(profiles))
	for i, p := range profiles {
		modelProfiles[i] = models.Profile{
			ID:        p.ID,
			Name:      p.Name,
			CreatedAt: p.CreatedAt,
			UpdatedAt: p.UpdatedAt,
		}
	}

	return modelProfiles, nil
}

// CreateProfile makes a new profile with validation
func (s *ProfileService) CreateProfile(ctx context.Context, profile models.Profile) (models.Profile, error) {
	// basic validation - name can't be empty
	if strings.TrimSpace(profile.Name) == "" {
		return models.Profile{}, errors.New("profile name cannot be empty")
	}

	// generate UUID if not provided
	if profile.ID == uuid.Nil {
		profile.ID = uuid.New()
	}

	// let database handle the creation
	createdProfile, err := s.DB.CreateProfile(ctx, database.CreateProfileParams{
		ID:   profile.ID,
		Name: profile.Name,
	})
	if err != nil {
		log.Printf("Error creating profile: %v", err)
		return models.Profile{}, fmt.Errorf("failed to create profile: %w", err)
	}

	// convert back to app model
	return models.Profile{
		ID:        createdProfile.ID,
		Name:      createdProfile.Name,
		CreatedAt: createdProfile.CreatedAt,
		UpdatedAt: createdProfile.UpdatedAt,
	}, nil
}

// UpdateProfileName updates profile name by user ID (changed from name-based to ID-based for safety)
func (s *ProfileService) UpdateProfileName(ctx context.Context, userID uuid.UUID, newName string) (models.Profile, error) {
	// validate inputs
	if userID == uuid.Nil {
		return models.Profile{}, errors.New("user ID cannot be empty")
	}

	if strings.TrimSpace(newName) == "" {
		return models.Profile{}, errors.New("new name cannot be empty")
	}

	// let database handle the update
	updatedProfile, err := s.DB.UpdateProfileByID(ctx, database.UpdateProfileByIDParams{
		ID:   userID,
		Name: newName,
	})
	if err != nil {
		log.Printf("Error updating profile by ID: %v", err)
		return models.Profile{}, fmt.Errorf("failed to update profile: %w", err)
	}

	// convert back to app model
	return models.Profile{
		ID:        updatedProfile.ID,
		Name:      updatedProfile.Name,
		CreatedAt: updatedProfile.CreatedAt,
		UpdatedAt: updatedProfile.UpdatedAt,
	}, nil
}

// GetProfileByID retrieves a profile by its ID
func (s *ProfileService) GetProfileByID(ctx context.Context, id uuid.UUID) (models.Profile, error) {
	// let database fetch the profile by ID
	dbProfile, err := s.DB.GetProfileById(ctx, id)
	if err != nil {
		log.Printf("Error retrieving profile by ID: %v", err)
		return models.Profile{}, fmt.Errorf("failed to get profile by ID: %w", err)
	}

	// convert back to app model
	return models.Profile{
		ID:        dbProfile.ID,
		Name:      dbProfile.Name,
		CreatedAt: dbProfile.CreatedAt,
		UpdatedAt: dbProfile.UpdatedAt,
	}, nil
}

// DeleteProfileByID deletes a profile by user ID (safer than name-based deletion)
func (s *ProfileService) DeleteProfileByID(ctx context.Context, userID uuid.UUID) error {
	// validate input
	if userID == uuid.Nil {
		return errors.New("user ID cannot be empty")
	}

	// let database handle the deletion
	if err := s.DB.DeleteProfile(ctx, userID); err != nil {
		log.Printf("Error deleting profile by ID: %v", err)
		return fmt.Errorf("failed to delete profile: %w", err)
	}

	return nil
}
