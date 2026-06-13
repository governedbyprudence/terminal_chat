package components

import (
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

var roomMenuItem = lipgloss.NewStyle().
	Foreground(Text).
	Padding(0, 2)

var roomSelectedMenuItem = lipgloss.NewStyle().
	Background(Accent).
	Foreground(DeepBg).
	Bold(true).
	Padding(0, 2)

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
			b.WriteString(roomSelectedMenuItem.Render("  > " + room))
		} else {
			b.WriteString(roomMenuItem.Render("    " + room))
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
		return RedirectSelect, nil
	}
	return RedirectNone, nil
}

func (r *RoomList) Update(msg tea.Msg) tea.Cmd { return nil }
func (r *RoomList) Focus() tea.Cmd             { return nil }
func (r *RoomList) Blur()                      {}
func (r *RoomList) Focused() bool              { return false }
