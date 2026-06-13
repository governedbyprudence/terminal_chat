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

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)
type Model struct {
	Width, Height int
	Screen        states.AppScreen
}

func (m Model) Init() tea.Cmd {
	f, err := os.OpenFile("debug.log",
		os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
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
		for raw := range sub.Channel() {
			var msg states.Message
			if err := json.Unmarshal([]byte(raw), &msg); err == nil {
				return msg
			}
		}
		return nil
	}
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.Width = msg.Width
		m.Height = msg.Height
		return m, nil

	case tea.KeyMsg:
		if msg.String() == "q" {
			return m, tea.Quit
		}
		s, cmd := m.Screen.Update(msg)
		m.Screen = s
		return m, cmd

	case states.Message:
		s, _ := m.Screen.Update(msg)
		m.Screen = s
		return m, m.listenForMessages()
	}

	return m, nil
}

func (m Model) View() string {
	minW := 50
	if m.Width < minW {
		return lipgloss.NewStyle().
			Foreground(styles.Accent).
			Bold(true).
			Padding(2, 2).
			Render(fmt.Sprintf("Terminal too narrow (%d cols).\nMinimum width: %d columns.", m.Width, minW))
	}
	return m.Screen.Render(m.Width, m.Height)
}

func newModel() Model {
	return Model{
		Screen: ui.NewRoomListScreen([]string{"Room 1", "Room 2", "Room 3"}),
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
