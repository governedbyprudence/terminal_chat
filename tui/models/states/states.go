package states

import tea "github.com/charmbracelet/bubbletea"

type Message struct {
	RoomName string `json:"room_name"`
	Text     string `json:"message"`
	UserID   string `json:"user_id"`
}

type PresenceEvent struct {
	UserID   string `json:"user_id"`
	RoomName string `json:"room_name"`
	Action   string `json:"action"` // "join" or "leave"
}

type AppScreen interface {
	Render(width, height int) string
	Update(msg tea.Msg) (AppScreen, tea.Cmd)
}
