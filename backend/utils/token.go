package utils

import (
	"crypto/rand"
	"encoding/hex"
)

// GenerateToken generates a cryptographically secure random token
func GenerateToken() string {
	bytes := make([]byte, 32)
	if _, err := rand.Read(bytes); err != nil {
		panic(err)
	}
	return hex.EncodeToString(bytes)
}
