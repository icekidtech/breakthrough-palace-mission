package utils

import (
	"time"

	"github.com/google/uuid"
)

// ParseUUID converts string to UUID, returns zero UUID if invalid
func ParseUUID(s string) uuid.UUID {
	u, err := uuid.Parse(s)
	if err != nil {
		return uuid.Nil
	}
	return u
}

// TimeNowAddHours returns current time plus hours
func TimeNowAddHours(hours int) time.Time {
	return time.Now().Add(time.Duration(hours) * time.Hour)
}
