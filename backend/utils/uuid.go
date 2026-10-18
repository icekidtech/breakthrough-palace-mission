package utils

import (
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
