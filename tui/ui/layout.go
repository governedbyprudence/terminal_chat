package ui

import (
	"terminal_chat/tui/components"

	"github.com/charmbracelet/lipgloss"
)

var outerFrame = lipgloss.NewStyle().
	Border(lipgloss.DoubleBorder(), true).
	BorderForeground(components.Accent).
	Padding(0, 2)

func outerFrameHSize() int {
	return outerFrame.GetHorizontalFrameSize()
}

func outerFrameVSize() int {
	return outerFrame.GetVerticalFrameSize()
}

