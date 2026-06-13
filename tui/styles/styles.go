package styles

import (
	"github.com/charmbracelet/bubbles/textinput"
	"github.com/charmbracelet/lipgloss"
)

// 5-color brutalist palette — every color has a functional semantic purpose.
var (
	DeepBg  = lipgloss.Color("#0D0D0D")
	Surface = lipgloss.Color("#1A1A1A")
	Text    = lipgloss.Color("#E8E8E8")
	DimText = lipgloss.Color("#6B6B6B")
	Accent  = lipgloss.Color("#FF6A00")
)

// OuterFrame — single accent double border.
// No fixed width — content determines size, so inner alignment isn't overridden.
var OuterFrame = lipgloss.NewStyle().
	Border(lipgloss.DoubleBorder(), true).
	BorderForeground(Accent).
	Padding(0, 2)

const OuterFrameHSize = 1 + 1 + 2 + 2 // border (║║) + padding (0+2+2)

// TitleBar — accent-filled header.
var TitleBar = lipgloss.NewStyle().
	Background(Accent).
	Foreground(DeepBg).
	Bold(true).
	Padding(0, 3)

// MenuItem — unselected room row.
var MenuItem = lipgloss.NewStyle().
	Foreground(Text).
	Padding(0, 2)

// SelectedMenuItem — selected room with accent background.
var SelectedMenuItem = lipgloss.NewStyle().
	Background(Accent).
	Foreground(DeepBg).
	Bold(true).
	Padding(0, 2)

// FooterText — dim navigation hint.
var FooterText = lipgloss.NewStyle().
	Foreground(DimText).
	Padding(0, 2)

// InputPrompt — accent prompt character.
var InputPrompt = lipgloss.NewStyle().
	Foreground(Accent).
	Bold(true)

// InputWrapper — top-bordered section separator above the input.
// Top border only — no side borders, so it doesn't affect horizontal width.
var InputWrapper = lipgloss.NewStyle().
	Border(lipgloss.NormalBorder(), true, false, false, false).
	BorderForeground(Accent).
	Padding(0, 2)

// SentMessage — own messages: bold weight (no padding — alignment handles it).
var SentMessage = lipgloss.NewStyle().
	Foreground(Text).
	Bold(true)

// ReceivedMessage — others' messages: normal weight (no padding).
var ReceivedMessage = lipgloss.NewStyle().
	Foreground(Text)

// SenderLabel — dim sender tag.
var SenderLabel = lipgloss.NewStyle().
	Foreground(DimText)

// MessageTimestamp — dim italic label.
var MessageTimestamp = lipgloss.NewStyle().
	Foreground(DimText).
	Italic(true)

// EmptyState — dim placeholder.
var EmptyState = lipgloss.NewStyle().
	Foreground(DimText).
	Italic(true).
	Padding(1, 2)

// TextInputStyle — input text style.
// Padding is intentionally excluded: textinput calls TextStyle.Inline(true)
// internally, which silently strips padding/borders/margins.
var TextInputStyle = lipgloss.NewStyle().
	Foreground(Text)

func NewTextInput(isPassword bool) textinput.Model {
	ti := textinput.New()
	ti.Prompt = ""
	ti.Focus()
	ti.TextStyle = TextInputStyle
	if isPassword {
		ti.EchoMode = textinput.EchoPassword
	}
	return ti
}
