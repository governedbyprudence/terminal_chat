package components

import (
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

var activeTab = lipgloss.NewStyle().
	Background(Accent).
	Foreground(DeepBg).
	Bold(true).
	Padding(0, 1)

var inactiveTab = lipgloss.NewStyle().
	Background(Surface).
	Foreground(Text).
	Padding(0, 1)

type RoomTabs struct {
	Rooms       []string
	ActiveIndex int
}

func NewRoomTabs(rooms []string) *RoomTabs {
	return &RoomTabs{Rooms: rooms}
}

func (t *RoomTabs) ActiveRoom() string {
	if t.ActiveIndex < 0 || t.ActiveIndex >= len(t.Rooms) {
		return ""
	}
	return t.Rooms[t.ActiveIndex]
}

func (t *RoomTabs) SelectNext() {
	if len(t.Rooms) == 0 {
		return
	}
	t.ActiveIndex = (t.ActiveIndex + 1) % len(t.Rooms)
}

func (t *RoomTabs) SelectPrev() {
	if len(t.Rooms) == 0 {
		return
	}
	t.ActiveIndex--
	if t.ActiveIndex < 0 {
		t.ActiveIndex = len(t.Rooms) - 1
	}
}

func (t *RoomTabs) SelectIndex(i int) {
	if i < 0 || i >= len(t.Rooms) {
		return
	}
	t.ActiveIndex = i
}

func truncate(s string, maxWidth int) string {
	if maxWidth < 1 {
		return ""
	}
	w := lipgloss.Width(s)
	if w <= maxWidth {
		return s
	}
	var b strings.Builder
	var cw int
	for _, r := range s {
		rw := lipgloss.Width(string(r))
		if cw+rw > maxWidth-1 {
			b.WriteRune('…')
			break
		}
		b.WriteRune(r)
		cw += rw
	}
	return b.String()
}

func (t *RoomTabs) Render(width int) string {
	n := len(t.Rooms)
	if n == 0 {
		return strings.Repeat(" ", width)
	}

	gap := 2
	totalGaps := (n - 1) * gap
	overheadPerTab := 6
	availableForNames := width - totalGaps - n*overheadPerTab
	nameWidth := availableForNames / n
	if nameWidth < 1 {
		nameWidth = 1
	}

	var rendered []string
	for i, room := range t.Rooms {
		displayName := truncate(room, nameWidth)
		var tab string
		if i == t.ActiveIndex {
			tab = activeTab.Render("╭ " + displayName + " ╮")
		} else {
			tab = inactiveTab.Render("╭ " + displayName + " ╮")
		}
		rendered = append(rendered, tab)
	}

	result := lipgloss.JoinHorizontal(lipgloss.Top, rendered...)
	resultWidth := lipgloss.Width(result)
	if resultWidth < width {
		result += strings.Repeat(" ", width-resultWidth)
	}
	return result
}

func (t *RoomTabs) OnKeyPress(key string) (int, tea.Cmd) {
	return 0, nil
}

func (t *RoomTabs) Update(msg tea.Msg) tea.Cmd { return nil }
func (t *RoomTabs) Focus() tea.Cmd             { return nil }
func (t *RoomTabs) Blur()                      {}
func (t *RoomTabs) Focused() bool              { return false }
