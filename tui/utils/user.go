package utils

import (
	"crypto/rand"
	"fmt"
	"sync"
	"time"
)

var (
	userID string
	once   sync.Once
)

func CurrentUserID() string {
	once.Do(func() {
		b := make([]byte, 8)
		_, err := rand.Read(b)
		if err != nil {
			userID = fmt.Sprintf("user_%d", time.Now().UnixNano())
		} else {
			userID = fmt.Sprintf("user_%x", b)
		}
	})
	return userID
}
