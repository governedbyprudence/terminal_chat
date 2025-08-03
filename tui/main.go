package main

import (
	"fmt"
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
	return nil
}

// Update Method
func (m Model) Update(msg bubbletea.Msg) (tea.Model, bubbletea.Cmd) {
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
		State: initState,
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
