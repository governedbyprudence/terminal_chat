package main

import (
	"fmt"
	"log"
	"os"
	"terminal_chat/tui/models/states"
	"terminal_chat/tui/ui"

	bubbletea "github.com/charmbracelet/bubbletea"
	tea "github.com/charmbracelet/bubbletea"
)

type Model struct {
	Screen states.AppScreen
}

// Init Method
func (m Model) Init() bubbletea.Cmd {
	f, err := os.OpenFile("debug.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Fatal(err)
	}
	log.SetOutput(f)
	return nil
}

// Update Method
func (m Model) Update(msg bubbletea.Msg) (tea.Model, bubbletea.Cmd) {

	switch msg := msg.(type) {
	case tea.KeyMsg:
		log.Println(msg.String())

		redirect := m.Screen.OnPressKey(msg.String())

		if redirect == 1 {
			// TODO: Enter Chat Screen

		}

		switch msg.String() {
		case "q":
			return m, tea.Quit
		}
	}

	return m, nil
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
	}

	p := tea.NewProgram(initModel, tea.WithAltScreen())

	if _, err := p.Run(); err != nil {
		fmt.Printf("There's been an error: %v", err)
		os.Exit(1)
	}
}
