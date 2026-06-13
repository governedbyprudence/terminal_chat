package components

import "github.com/charmbracelet/lipgloss"

var (
	DeepBg  = lipgloss.Color("#0D0D0D")
	Surface = lipgloss.Color("#1A1A1A")
	Text    = lipgloss.Color("#E8E8E8")
	DimText = lipgloss.Color("#6B6B6B")
	Accent  = lipgloss.Color("#FF6A00")
)

var (
	MenuItem         = lipgloss.NewStyle().Foreground(Text).Padding(0, 2)
	SelectedMenuItem = lipgloss.NewStyle().Background(Accent).Foreground(DeepBg).Bold(true).Padding(0, 2)
)
