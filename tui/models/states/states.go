package states

import tea "github.com/charmbracelet/bubbletea"

type Message struct {
	RoomName string `json:"room_name"`
	Text     string `json:"message"`
	UserID   string `json:"user_id"`
}

type AppScreen interface {
	Render(width, height int) string
	Update(msg tea.Msg) (AppScreen, tea.Cmd)
}
