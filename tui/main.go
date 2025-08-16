package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"terminal_chat/tui/messageStore"
	"terminal_chat/tui/models/states"
	"terminal_chat/tui/styles"
	"terminal_chat/tui/ui"
	"terminal_chat/tui/utils"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/redis/go-redis/v9"
)

var ctx = context.Background()

type StoreMessage struct {
	RoomName string `json:"room_name"`
	Message  string `json:"message"`
	UserID   string `json:"user_id"`
}

// CurrentUserID returns the current user's ID
func CurrentUserID() string {
	return utils.CurrentUserID()
}

type Model struct {
	Redis   *redis.Client
	program *tea.Program
	Screen  states.AppScreen
	Screens map[string]states.AppScreen
}

// Init Method
func (m Model) Init() tea.Cmd {
	m.Redis = messageStore.GetClient()
	f, err := os.OpenFile("debug.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Fatal(err)
	}
	log.SetOutput(f)
	// Start the Redis subscription
	return m.listenForMessages()
}

// listenForMessages starts a Redis subscription and returns a command that processes messages
func (m Model) listenForMessages() tea.Cmd {
	return func() tea.Msg {
		sub := m.Redis.Subscribe(context.Background(), "chat")
		channel := sub.Channel()

		for msg := range channel {
			var message StoreMessage
			if err := json.Unmarshal([]byte(msg.Payload), &message); err == nil {
				return message
			}
		}
		return nil
	}
}

// Update Method
func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

	switch msg := msg.(type) {
	// Handle the initial command from Init
	case tea.Cmd:
		return m, msg

	case StoreMessage:
		// Check if the current screen is ChatRoomScreen and the room name matches
		if chatScreen, ok := m.Screen.(ui.ChatRoomScreen); ok {
			if chatScreen.State.RoomName == msg.RoomName {
				// Add the message to the messages slice with the user ID
				chatScreen.UpdateMessage(msg.Message, msg.UserID)
				m.Screen = chatScreen
			}
		}
		// Continue listening for more messages
		return m, m.listenForMessages()
	case tea.KeyMsg:
		log.Println(msg.String())

		redirect, cmd := m.Screen.OnPressKey(msg.String())

		if cmd != nil {
			return m, cmd
		}

		if redirect == 1 {
			// If current screen is RoomListScreen
			if _, ok := m.Screen.(ui.RoomListScreen); ok {
				// Add the current screen to the screens map
				m.Screens["roomlist"] = m.Screen

				// Create a new text input
				ti := styles.NewTextInput(false)

				// Update the current screen
				m.Screen = ui.ChatRoomScreen{
					State: &states.ChatState{
						Redis:     m.Redis,
						TextInput: ti,
						RoomName:  m.Screen.(ui.RoomListScreen).State.Rooms[m.Screen.(ui.RoomListScreen).State.SelectedRoom],
					},
				}

				return m, textinput.Blink
			} else if _, ok := m.Screen.(ui.ChatRoomScreen); ok {
				// Add the current screen to the screens map
				m.Screens["chatroom"] = m.Screen

				// Update the current screen
				m.Screen = m.Screens["roomlist"]
			}
		}

		if msg.String() == "q" {
			return m, tea.Quit
		}

		// Handle text input updates
		if _, ok := m.Screen.(ui.ChatRoomScreen); ok {
			input := m.Screen.(ui.ChatRoomScreen).State.TextInput
			input, cmd = input.Update(msg)
			m.Screen.(ui.ChatRoomScreen).State.TextInput = input
			return m, cmd
		}
	}

	return m, cmd
}

// View Method
func (m Model) View() string {
	return m.Screen.Render()
}

func main() {
	initState := states.RoomListState{
		Rooms: []string{"Room 1", "Room 2", "Room 3"},
	}

	initScreen := ui.RoomListScreen{
		State: &initState,
	}

	initModel := Model{
		Screen: initScreen,
		Screens: map[string]states.AppScreen{
			"roomlist": initScreen,
			"chatroom": ui.ChatRoomScreen{
				State: &states.ChatState{
					RoomName: "Room 1",
				},
			},
		},
		Redis: messageStore.GetClient(),
	}

	p := tea.NewProgram(initModel, tea.WithAltScreen())

	initModel.program = p

	if _, err := p.Run(); err != nil {
		fmt.Printf("There's been an error: %v", err)
		os.Exit(1)
	}
}
