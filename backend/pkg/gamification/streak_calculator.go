package gamification

import (
	"database/sql"
	"time"
)

func CalculateNewStreak(lastActiveDate sql.NullTime, currentStreak int) (newStreak int, shouldUpdate bool) {
	if !lastActiveDate.Valid {
		// first time ever, start the streak
		return 1, true
	}

	now := time.Now()
	lastActive := lastActiveDate.Time

	// get the dates without time components for comparison
	todayStart := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())
	lastActiveStart := time.Date(lastActive.Year(), lastActive.Month(), lastActive.Day(), 0, 0, 0, 0, lastActive.Location())

	daysSince := int(todayStart.Sub(lastActiveStart).Hours() / 24)

	if daysSince == 0 {
		// same day, don't update streak
		return currentStreak, false
	} else if daysSince == 1 {
		// consecutive day, increment streak
		return currentStreak + 1, true
	} else {
		// missed a day or more, reset streak
		return 1, true
	}
}
