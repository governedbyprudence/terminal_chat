package ui

import (
	"terminal_chat/tui/models/states"
	"terminal_chat/tui/styles"

	tea "github.com/charmbracelet/bubbletea"
	lipgloss "github.com/charmbracelet/lipgloss"
)

type RoomListScreen struct {
	State *states.RoomListState
}

func (s RoomListScreen) Render() string {
	title := styles.TitleStyle.Render("Room List")

	body := ""

	for i, room := range s.State.Rooms {
		if i == s.State.SelectedRoom {
			body += styles.SelectedMenuItemStyle.Render(room) + "\n"
		} else {
			body += styles.MenuItemStyle.Render(room) + "\n"
		}
	}

	footer := styles.FooterStyle.Render("Use ↑ ↓ to navigate the menu.\n\nPress `q` to quit.")

	return lipgloss.JoinVertical(lipgloss.Top, title, body, footer)
}

func (s RoomListScreen) OnPressKey(key string) (int, tea.Cmd) {

	switch key {
	case "down":
		if s.State.SelectedRoom < len(s.State.Rooms)-1 {
			s.State.SelectedRoom++
		}
	case "up":
		if s.State.SelectedRoom > 0 {
			s.State.SelectedRoom--
		}
	case "enter":
		return 1, nil // Redirect to chat screen
	}

	return 0, nil // No redirect
}
