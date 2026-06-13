package components

import (
	"strings"

	"github.com/charmbracelet/lipgloss"
)

var userPanelTitle = lipgloss.NewStyle().
	Background(Accent).
	Foreground(DeepBg).
	Bold(true).
	Padding(0, 1)

var userPanelStyle = lipgloss.NewStyle().
	Border(lipgloss.NormalBorder(), false, false, false, true).
	BorderForeground(Accent).
	Padding(0, 1)

const (
	RightPanelWidth      = 22
	PanelGap             = 2
	userPanelInnerAdjust = 5
	userPanelBorderAdj   = 3
)

type UserPanel struct {
	UserList *UserList
}

func NewUserPanel() *UserPanel {
	return &UserPanel{UserList: NewUserList()}
}

func (p *UserPanel) RenderContent() string {
	var b strings.Builder
	b.WriteString(userPanelTitle.Render("ONLINE"))
	b.WriteString("\n")
	if p.UserList.Count() == 0 {
		b.WriteString("nobody here")
		b.WriteString("\n")
	} else {
		b.WriteString(p.UserList.Render(RightPanelWidth - userPanelInnerAdjust))
	}
	return strings.TrimSuffix(b.String(), "\n")
}

func (p *UserPanel) RenderFull(targetLines int) string {
	content := p.RenderContent()
	contentLines := strings.Count(content, "\n") + 1

	var gapBuf strings.Builder
	for i := 0; i < targetLines-1; i++ {
		gapBuf.WriteString("  \n")
	}
	gapBuf.WriteString("  ")

	var padBuf strings.Builder
	padBuf.WriteString(content)
	for i := contentLines; i < targetLines; i++ {
		padBuf.WriteString("\n")
	}
	paddedContent := userPanelStyle.Width(RightPanelWidth - userPanelBorderAdj).Render(padBuf.String())

	return lipgloss.JoinHorizontal(lipgloss.Top, gapBuf.String(), paddedContent)
}
