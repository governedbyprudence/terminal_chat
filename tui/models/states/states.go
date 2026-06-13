package states

import (
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
)

type Message struct {
	RoomName string `json:"room_name"`
	Text     string `json:"message"`
	UserID   string `json:"user_id"`
}

type RoomListState struct {
	Rooms        []string
	SelectedRoom int
}

type ChatState struct {
	RoomName  string
	Messages  []Message
	TextInput textinput.Model
}

type AppScreen interface {
	Render(width, height int) string
	OnPressKey(key string) (int, tea.Cmd)
}
