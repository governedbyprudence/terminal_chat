package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"os"

	"terminal_chat/tui/messageStore"
	"terminal_chat/tui/models/states"
	"terminal_chat/tui/styles"
	"terminal_chat/tui/ui"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
)

type Model struct {
	Screen  states.AppScreen
	Screens map[string]states.AppScreen
}

func (m Model) Init() tea.Cmd {
	f, err := os.OpenFile("debug.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Fatal(err)
	}
	log.SetOutput(f)
	return m.listenForMessages()
}

func (m Model) listenForMessages() tea.Cmd {
	return func() tea.Msg {
		sub := messageStore.DefaultBroker().Subscribe("chat")
		defer sub.Close()
		for msg := range sub.Channel() {
			var m states.Message
			if err := json.Unmarshal([]byte(msg), &m); err == nil {
				return m
			}
		}
		return nil
	}
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.Cmd:
		return m, msg

	case states.Message:
		chat, ok := m.Screen.(ui.ChatRoomScreen)
		if ok && chat.State.RoomName == msg.RoomName {
			chat.UpdateMessage(msg)
			m.Screen = chat
		}
		return m, m.listenForMessages()

	case tea.KeyMsg:
		log.Println(msg.String())

		redirect, cmd := m.Screen.OnPressKey(msg.String())

		if cmd != nil {
			return m, cmd
		}

		if redirect == 1 {
			if _, ok := m.Screen.(ui.RoomListScreen); ok {
				m.Screens["roomlist"] = m.Screen

				ti := styles.NewTextInput(false)

				m.Screen = ui.ChatRoomScreen{
					State: &states.ChatState{
						TextInput: ti,
						RoomName:  m.Screen.(ui.RoomListScreen).State.Rooms[m.Screen.(ui.RoomListScreen).State.SelectedRoom],
					},
				}

				return m, textinput.Blink
			} else if _, ok := m.Screen.(ui.ChatRoomScreen); ok {
				m.Screens["chatroom"] = m.Screen
				m.Screen = m.Screens["roomlist"]
			}
		}

		if msg.String() == "q" {
			return m, tea.Quit
		}

		if _, ok := m.Screen.(ui.ChatRoomScreen); ok {
			input := m.Screen.(ui.ChatRoomScreen).State.TextInput
			input, cmd = input.Update(msg)
			m.Screen.(ui.ChatRoomScreen).State.TextInput = input
			return m, cmd
		}
	}

	return m, nil
}

func (m Model) View() string {
	return m.Screen.Render()
}

func newModel() Model {
	initState := states.RoomListState{
		Rooms: []string{"Room 1", "Room 2", "Room 3"},
	}

	initScreen := ui.RoomListScreen{
		State: &initState,
	}

	return Model{
		Screen: initScreen,
		Screens: map[string]states.AppScreen{
			"roomlist": initScreen,
			"chatroom": ui.ChatRoomScreen{
				State: &states.ChatState{
					RoomName: "Room 1",
				},
			},
		},
	}
}

func runLocal() {
	initModel := newModel()
	p := tea.NewProgram(initModel, tea.WithAltScreen())
	if _, err := p.Run(); err != nil {
		fmt.Printf("There's been an error: %v", err)
		os.Exit(1)
	}
}

func main() {
	serverAddr := flag.String("server", "", "Start TCP server on given address (e.g. :9000)")
	connectAddr := flag.String("connect", "", "Connect to TCP server (e.g. localhost:9000)")
	flag.Parse()

	if *serverAddr != "" {
		if err := messageStore.StartServer(*serverAddr); err != nil {
			log.Fatal(err)
		}
		fmt.Println("TCP server started on", *serverAddr)
		if err := messageStore.DefaultBroker().Connect(*serverAddr); err != nil {
			log.Fatal(err)
		}
		fmt.Println("Connected to local TCP server")
	} else if *connectAddr != "" {
		if err := messageStore.DefaultBroker().Connect(*connectAddr); err != nil {
			log.Fatal(err)
		}
		fmt.Println("Connected to TCP server at", *connectAddr)
	}

	runLocal()
}
