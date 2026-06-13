package ui_test

import (
	"fmt"
	"strings"
	"testing"

	"terminal_chat/tui/messageStore"
	"terminal_chat/tui/models"
	"terminal_chat/tui/ui"

	tea "github.com/charmbracelet/bubbletea"
)

func TestChatRoomScreen_Render_LineCount(t *testing.T) {
	broker := messageStore.New()
	backRef := ui.NewRoomListScreen(broker, []string{"back"})
	screen := ui.NewChatRoomScreen(broker, backRef, "test",
		models.Message{UserID: "alice", Text: "hi"},
		models.Message{UserID: "bob", Text: "hello"},
	)

	width, height := 74, 24
	out := screen.Render(width, height)
	lines := strings.Split(out, "\n")

	if len(lines) != height {
		t.Errorf("expected %d lines, got %d", height, len(lines))
	}
	if !strings.Contains(out, "test") {
		t.Error("expected room name in output")
	}
	if !strings.Contains(out, "╔") || !strings.Contains(out, "╚") {
		t.Error("expected border characters")
	}
}

func TestChatRoomScreen_Render_Small(t *testing.T) {
	broker := messageStore.New()
	backRef := ui.NewRoomListScreen(broker, []string{"back"})
	screen := ui.NewChatRoomScreen(broker, backRef, "test")

	width, height := 74, 10
	out := screen.Render(width, height)
	lines := strings.Split(out, "\n")

	if len(lines) != height {
		t.Errorf("expected %d lines, got %d", height, len(lines))
	}
	if !strings.Contains(out, "╔") || !strings.Contains(out, "╚") {
		t.Error("expected border characters")
	}
}

func TestChatRoomScreen_Render_NoTrailingNL(t *testing.T) {
	broker := messageStore.New()
	backRef := ui.NewRoomListScreen(broker, []string{"back"})
	screen := ui.NewChatRoomScreen(broker, backRef, "test")

	out := screen.Render(74, 24)
	if strings.HasSuffix(out, "\n") {
		t.Error("screen Render should not have trailing newline")
	}
}

func TestChatRoomScreen_Update_Message(t *testing.T) {
	broker := messageStore.New()
	backRef := ui.NewRoomListScreen(broker, []string{"back"})
	screen := ui.NewChatRoomScreen(broker, backRef, "test",
		models.Message{UserID: "alice", Text: "hi"},
	)

	msg := models.Message{RoomName: "test", UserID: "bob", Text: "hello"}
	next, _ := screen.Update(msg)
	if next == nil {
		t.Fatal("Update returned nil screen")
	}

	out := next.Render(74, 24)
	if !strings.Contains(out, "hello") {
		t.Error("expected new message in rendered output")
	}
}

func TestChatRoomScreen_Update_Escape(t *testing.T) {
	broker := messageStore.New()
	backRef := ui.NewRoomListScreen(broker, []string{"back"})
	screen := ui.NewChatRoomScreen(broker, backRef, "test")

	next, _ := screen.Update(tea.KeyMsg{Type: tea.KeyEsc})
	if next != backRef {
		t.Error("esc should return backRef")
	}
}

func TestChatRoomScreen_Update_PgUpDown(t *testing.T) {
	broker := messageStore.New()
	backRef := ui.NewRoomListScreen(broker, []string{"back"})
	screen := ui.NewChatRoomScreen(broker, backRef, "test")

	for i := 0; i < 15; i++ {
		screen.Update(models.Message{
			RoomName: "test", UserID: "u", Text: fmt.Sprintf("msg %d", i+1),
		})
	}

	outBefore := screen.Render(74, 24)

	next, _ := screen.Update(tea.KeyMsg{Type: tea.KeyPgUp})
	outAfter := next.Render(74, 24)
	if outBefore == outAfter {
		t.Error("pgup should change visible content")
	}

	next2, _ := next.Update(tea.KeyMsg{Type: tea.KeyPgDown})
	outAfter2 := next2.Render(74, 24)
	if outBefore != outAfter2 {
		t.Error("pgdown after pgup should restore original")
	}
}

func TestRoomListScreen_Render(t *testing.T) {
	broker := messageStore.New()
	screen := ui.NewRoomListScreen(broker, []string{"Room A", "Room B"})

	out := screen.Render(74, 24)
	if !strings.Contains(out, "ROOMS") {
		t.Error("expected ROOMS title")
	}
	if !strings.Contains(out, "Room A") || !strings.Contains(out, "Room B") {
		t.Error("expected room names")
	}
}

func TestRoomListScreen_Update_Enter(t *testing.T) {
	broker := messageStore.New()
	screen := ui.NewRoomListScreen(broker, []string{"Room A"})

	next, cmd := screen.Update(tea.KeyMsg{Type: tea.KeyEnter})
	if next == nil {
		t.Fatal("Update returned nil")
	}
	if cmd == nil {
		t.Error("enter should return a command (textinput.Blink)")
	}
	out := next.Render(74, 24)
	if !strings.Contains(out, "Room A") {
		t.Error("expected room name in chat screen")
	}
}

func TestScreen_Message_RouteNoPanic(t *testing.T) {
	broker := messageStore.New()
	screen := ui.NewRoomListScreen(broker, []string{"Room A"})
	msg := models.Message{RoomName: "nowhere", UserID: "x", Text: "y"}
	next, _ := screen.Update(msg)
	if next == nil {
		t.Fatal("Update returned nil")
	}
}

func TestChatRoomScreen_Update_EnterPublishes(t *testing.T) {
	broker := messageStore.New()
	sub := broker.Subscribe("chat")
	defer sub.Close()

	backRef := ui.NewRoomListScreen(broker, []string{"back"})
	screen := ui.NewChatRoomScreen(broker, backRef, "test")

	screen.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune("test msg")})
	next, _ := screen.Update(tea.KeyMsg{Type: tea.KeyEnter})
	if next == nil {
		t.Fatal("Update returned nil")
	}

	select {
	case payload := <-sub.Channel():
		if !strings.Contains(payload, "test msg") {
			t.Errorf("expected message text in published payload, got %q", payload)
		}
	default:
		t.Error("enter should publish message to channel")
	}
}
