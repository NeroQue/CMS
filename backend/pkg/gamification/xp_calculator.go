package gamification

import "math"

// XP per level constants
const (
	XPPerLevel   = 100 // Every 100 XP = 1 level
	GemsPerLevel = 1   // Award 1 gem per level up
	BaseXP       = 5   // Minimum XP for any content
)

// CalculateXPForContent dynamically calculates XP based on content characteristics
// Returns XP in increments of 5
func CalculateXPForContent(contentType string, duration int, size int64) int {
	var xp int

	switch contentType {
	case "video", "mp4", "avi", "mov", "mkv", "webm", "flv", "wmv":
		// 5 XP per minute of video
		minutes := float64(duration) / 60.0
		xp = int(math.Round(minutes * 5))

	case "pdf":
		// Estimate ~50KB per page, 10 XP per page
		estimatedPages := float64(size) / 50000.0
		xp = int(math.Round(estimatedPages * 10))

	case "text", "txt", "md", "markdown":
		// 5 XP per 1KB of text
		kb := float64(size) / 1000.0
		xp = int(math.Round(kb * 5))

	default:
		// Generic: 10 XP per 100KB
		xp = int(float64(size) / 100000.0 * 10)
	}

	// Round to nearest 5
	xp = int(math.Round(float64(xp)/5.0) * 5)

	// Ensure minimum XP
	if xp < BaseXP {
		xp = BaseXP
	}

	return xp
}

// CalculateLevelFromXP determines user level based on total XP
func CalculateLevelFromXP(totalXP int) int {
	return (totalXP / XPPerLevel) + 1
}

// GetXPForNextLevel returns XP needed to reach next level
func GetXPForNextLevel(currentXP int) int {
	currentLevel := CalculateLevelFromXP(currentXP)
	nextLevelXP := currentLevel * XPPerLevel
	return nextLevelXP - currentXP
}

// GetXPProgressInLevel returns XP progress within current level (0-99)
func GetXPProgressInLevel(totalXP int) int {
	return totalXP % XPPerLevel
}

// CalculateGemsForLevelUp checks if a level up occurred and returns gems to award
func CalculateGemsForLevelUp(oldXP, newXP int) int {
	oldLevel := CalculateLevelFromXP(oldXP)
	newLevel := CalculateLevelFromXP(newXP)

	levelsGained := newLevel - oldLevel
	if levelsGained > 0 {
		return levelsGained * GemsPerLevel
	}
	return 0
}
