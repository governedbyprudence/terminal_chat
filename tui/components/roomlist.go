package components

import (
	"strings"

	"terminal_chat/tui/styles"

	tea "github.com/charmbracelet/bubbletea"
)

type RoomList struct {
	Rooms        []string
	SelectedRoom int
}

func NewRoomList(rooms []string) *RoomList {
	return &RoomList{Rooms: rooms}
}

func (r *RoomList) Render(_ int) string {
	var b strings.Builder
	for i, room := range r.Rooms {
		if i == r.SelectedRoom {
			b.WriteString(styles.SelectedMenuItem.Render("  > " + room))
		} else {
			b.WriteString(styles.MenuItem.Render("    " + room))
		}
		b.WriteString("\n")
	}
	return b.String()
}

func (r *RoomList) OnKeyPress(key string) (int, tea.Cmd) {
	switch key {
	case "down":
		if r.SelectedRoom < len(r.Rooms)-1 {
			r.SelectedRoom++
		}
	case "up":
		if r.SelectedRoom > 0 {
			r.SelectedRoom--
		}
	case "enter":
		return 1, nil
	}
	return 0, nil
}

func (r *RoomList) Focus() tea.Cmd { return nil }
func (r *RoomList) Blur()          {}
func (r *RoomList) Focused() bool  { return false }
