package gamification

import (
	"database/sql"
	"time"
)

func CalculateNewStreak(lastActiveDate sql.NullTime, currentStreak int) (newStreak int, shouldUpdate bool) {
	if !lastActiveDate.Valid {
		return 1, true
	}

	hoursSince := time.Since(lastActiveDate.Time).Hours()

	if hoursSince < 24 {
		return currentStreak, false
	} else if hoursSince < 48 {
		return currentStreak + 1, true
	} else {
		return 1, true
	}
}
