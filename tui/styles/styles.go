package styles

import "github.com/charmbracelet/lipgloss"

var titleStyle = lipgloss.NewStyle().
	Foreground(lipgloss.Color("205")).
	Bold(true).
	Padding(1, 2)

var sidebarStyle = lipgloss.NewStyle().
	Border(lipgloss.NormalBorder()).
	Width(20).
	Padding(1, 1).
	Margin(0, 1)

var mainStyle = lipgloss.NewStyle().
	Border(lipgloss.RoundedBorder()).
	Padding(1, 2).
	Width(50)

var menuItemStyle = lipgloss.NewStyle().
	Foreground(lipgloss.Color("241"))

var selectedMenuItemStyle = lipgloss.NewStyle().
	Foreground(lipgloss.Color("229")).
	Background(lipgloss.Color("57")).
	Bold(true)

var footerStyle = lipgloss.NewStyle().
	Foreground(lipgloss.Color("241")).
	Italic(true).
	MarginTop(1).
	Padding(0, 2)
