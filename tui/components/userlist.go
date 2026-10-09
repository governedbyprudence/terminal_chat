package components

import (
	"strings"

	"terminal_chat/tui/models"

	tea "github.com/charmbracelet/bubbletea"
)

type UserList struct {
	users []string
}

func NewUserList() *UserList {
	return &UserList{}
}

func (u *UserList) AddUser(id string) {
	for _, user := range u.users {
		if user == id {
			return
		}
	}
	u.users = append(u.users, id)
}

func (u *UserList) RemoveUser(id string) {
	for i, user := range u.users {
		if user == id {
			u.users = append(u.users[:i], u.users[i+1:]...)
			return
		}
	}
}

func (u *UserList) Count() int {
	return len(u.users)
}

func (u *UserList) Render(_ int) string {
	var b strings.Builder
	for _, user := range u.users {
		shortID := user
		if len(shortID) > 6 {
			shortID = shortID[:6]
		}
		if user == models.CurrentUserID() {
			b.WriteString(SelectedMenuItem.Render("  " + shortID))
		} else {
			b.WriteString(MenuItem.Render("  " + shortID))
		}
		b.WriteString("\n")
	}
	return b.String()
}

func (u *UserList) OnKeyPress(key string) (int, tea.Cmd) { return 0, nil }
func (u *UserList) Update(msg tea.Msg) tea.Cmd           { return nil }
func (u *UserList) Focus() tea.Cmd                       { return nil }
func (u *UserList) Blur()                                {}
func (u *UserList) Focused() bool                        { return false }
