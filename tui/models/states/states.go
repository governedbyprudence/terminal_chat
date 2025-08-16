package states

import (
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/redis/go-redis/v9"
)

type Messages struct {
	RoomName string
	Message  string
	UserID   string // ID of the user who sent the message
}

type AppState any

type RoomListState struct {
	Rooms        []string
	SelectedRoom int
}

type ChatState struct {
	RoomName  string
	Messages  []Messages
	isTyping  bool
	TextInput textinput.Model
	Redis     *redis.Client
}

type AppScreen interface {
	Render() string
	OnPressKey(key string) (int, tea.Cmd)
}
