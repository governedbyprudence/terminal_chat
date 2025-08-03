package styles

import "github.com/charmbracelet/lipgloss"

// Colors
var (
	PrimaryColor    = lipgloss.Color("205")
	SecondaryColor  = lipgloss.Color("24")
	HighlightColor  = lipgloss.Color("11")
	BackgroundColor = lipgloss.Color("235")
)

// Title style with centered text and improved appearance
var TitleStyle = lipgloss.NewStyle().
	Foreground(PrimaryColor).
	Bold(true).
	Padding(1, 2).
	Align(lipgloss.Center)

// Menu item styles
var MenuItemStyle = lipgloss.NewStyle().
	Foreground(SecondaryColor).
	Padding(0, 1)

// Sidebar with improved border and styling
var SidebarStyle = lipgloss.NewStyle().
	Border(lipgloss.NormalBorder(), true).
	BorderForeground(SecondaryColor).
	Width(25).
	Padding(1, 1).
	Margin(0, 1).
	Align(lipgloss.Center)

// Main content area with rounded borders and better styling
var MainStyle = lipgloss.NewStyle().
	Padding(1, 2).
	Align(lipgloss.Center)

var SelectedMenuItemStyle = lipgloss.NewStyle().
	Foreground(lipgloss.Color("229")).
	Background(lipgloss.Color("57")).
	Padding(0, 1).
	Bold(true)

var FooterStyle = lipgloss.NewStyle().
	Foreground(lipgloss.Color("241")).
	Italic(true).
	MarginTop(1).
	Padding(0, 2)
