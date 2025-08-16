package utils

import (
	"crypto/rand"
	"fmt"
	"time"
)

var currentUserID string

// CurrentUserID returns the current user's ID
func CurrentUserID() string {
	return currentUserID
}

func init() {
	// Generate a random user ID when the package is initialized
	currentUserID = generateRandomID()
}

func generateRandomID() string {
	b := make([]byte, 8)
	_, err := rand.Read(b)
	if err != nil {
		// Fallback to timestamp if random generation fails
		return fmt.Sprintf("user_%d", time.Now().UnixNano())
	}
	return fmt.Sprintf("user_%x", b)
}
