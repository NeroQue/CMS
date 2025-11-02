package services

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"time"

	"github.com/NeroQue/course-management-backend/internal/database"
	"github.com/NeroQue/course-management-backend/pkg/gamification"
	"github.com/google/uuid"
)

// GamificationService handles all XP, level, and gem logic
type GamificationService struct {
	DB *database.Queries
}

// NewGamificationService creates a new gamification service
func NewGamificationService(db *database.Queries) *GamificationService {
	return &GamificationService{DB: db}
}

// XPAwardResult contains information about XP awarded
type XPAwardResult struct {
	XPAwarded      int  `json:"xp_awarded"`
	GemsAwarded    int  `json:"gems_awarded"`
	LeveledUp      bool `json:"leveled_up"`
	NewLevel       int  `json:"new_level"`
	OldLevel       int  `json:"old_level"`
	TotalXP        int  `json:"total_xp"`
	AlreadyAwarded bool `json:"already_awarded"` // True if XP was previously awarded
}

// AwardXPForContent awards XP when a content item is completed
// Returns XPAwardResult with details about what was awarded
// Prevents double-awarding even after progress resets
func (s *GamificationService) AwardXPForContent(ctx context.Context, userID uuid.UUID, contentItemID uuid.UUID) (*XPAwardResult, error) {
	result := &XPAwardResult{}

	// Check if XP was already awarded for this content (prevents double-dipping)
	progressCheck, err := s.DB.CheckIfXPAwarded(ctx, database.CheckIfXPAwardedParams{
		UserID:        userID,
		ContentItemID: contentItemID,
	})

	if err == nil && progressCheck.XpAwarded {
		// XP was already awarded, don't award again
		log.Printf("XP already awarded for user %s, content %s (amount: %d)", userID, contentItemID, progressCheck.XpAmount)
		result.AlreadyAwarded = true
		result.XPAwarded = 0
		return result, nil
	}

	// Get content item details
	contentItem, err := s.DB.GetContentItemByID(ctx, contentItemID)
	if err != nil {
		return nil, fmt.Errorf("failed to get content item: %w", err)
	}

	// Calculate XP (use override if set, otherwise calculate)
	var xpToAward int
	if contentItem.XpValue.Valid && contentItem.XpValue.Int32 > 0 {
		xpToAward = int(contentItem.XpValue.Int32)
		log.Printf("Using manual XP override: %d XP for content %s", xpToAward, contentItemID)
	} else {
		xpToAward = gamification.CalculateXPForContent(
			contentItem.ContentType,
			int(contentItem.Duration.Int32),
			contentItem.Size.Int64,
		)
		log.Printf("Calculated XP: %d for content %s (type: %s, duration: %d, size: %d)",
			xpToAward, contentItemID, contentItem.ContentType, contentItem.Duration.Int32, contentItem.Size.Int64)
	}

	// Get current profile
	profile, err := s.DB.GetProfileById(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to get profile: %w", err)
	}

	oldXP := int(profile.Experience)
	oldLevel := gamification.CalculateLevelFromXP(oldXP)
	newXP := oldXP + xpToAward
	newLevel := gamification.CalculateLevelFromXP(newXP)

	// Check for level up and award gems
	gemsAwarded := gamification.CalculateGemsForLevelUp(oldXP, newXP)

	log.Printf("Awarding XP: user %s, old XP: %d (level %d), new XP: %d (level %d), gems: %d",
		userID, oldXP, oldLevel, newXP, newLevel, gemsAwarded)

	// Update profile with new XP, gems, and level
	_, err = s.DB.UpdateProfileGamification(ctx, database.UpdateProfileGamificationParams{
		ID:             userID,
		Experience:     int32(newXP),
		Gems:           profile.Gems + int32(gemsAwarded),
		Level:          int32(newLevel),
		LastActiveDate: sql.NullTime{Time: time.Now(), Valid: true},
	})

	if err != nil {
		return nil, fmt.Errorf("failed to update profile: %w", err)
	}

	// Build result
	result.XPAwarded = xpToAward
	result.GemsAwarded = gemsAwarded
	result.LeveledUp = newLevel > oldLevel
	result.NewLevel = newLevel
	result.OldLevel = oldLevel
	result.TotalXP = newXP
	result.AlreadyAwarded = false

	log.Printf("XP award complete: %+v", result)

	return result, nil
}

// GetLevelProgress returns detailed level progress for a profile
func (s *GamificationService) GetLevelProgress(ctx context.Context, userID uuid.UUID) (map[string]interface{}, error) {
	profile, err := s.DB.GetProfileById(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to get profile: %w", err)
	}

	totalXP := int(profile.Experience)
	currentLevel := gamification.CalculateLevelFromXP(totalXP)
	xpInLevel := gamification.GetXPProgressInLevel(totalXP)
	xpToNext := gamification.GetXPForNextLevel(totalXP)

	return map[string]interface{}{
		"total_xp":         totalXP,
		"current_level":    currentLevel,
		"xp_in_level":      xpInLevel,
		"xp_to_next_level": xpToNext,
		"gems":             profile.Gems,
		"progress_pct":     float64(xpInLevel) / float64(gamification.XPPerLevel) * 100,
	}, nil
}
