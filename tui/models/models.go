package models

import (
	"crypto/rand"
	"fmt"
	"sync"
	"time"

	tea "github.com/charmbracelet/bubbletea"
)

type Message struct {
	RoomName string `json:"room_name"`
	Text     string `json:"message"`
	UserID   string `json:"user_id"`
}

type PresenceEvent struct {
	UserID   string `json:"user_id"`
	RoomName string `json:"room_name"`
	Action   string `json:"action"`
}

type AppScreen interface {
	Render(width, height int) string
	Update(msg tea.Msg) (AppScreen, tea.Cmd)
}

var (
	userID  string
	uidOnce sync.Once
)

func CurrentUserID() string {
	uidOnce.Do(func() {
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
