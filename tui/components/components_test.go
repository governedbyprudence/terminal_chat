package components

import (
	"strings"
	"testing"

	"terminal_chat/tui/models"
)

func TestMessageList_Render_Empty(t *testing.T) {
	ml := NewMessageList()
	out := ml.Render(50, 6)
	lines := strings.Split(out, "\n")
	if len(lines) != 6 {
		t.Errorf("expected 6 lines, got %d; output=%q", len(lines), out)
	}
	if !strings.Contains(out, "no messages yet") {
		t.Error("expected empty state text")
	}
	if strings.HasSuffix(out, "\n") {
		t.Error("trailing newline should be stripped")
	}
}

func TestMessageList_Render_NonEmpty_NoTrailingNL(t *testing.T) {
	ml := NewMessageList()
	ml.AddMessage(models.Message{UserID: "alice", Text: "hi"})
	ml.AddMessage(models.Message{UserID: "bob", Text: "hello"})

	out := ml.Render(50, 8)
	lines := strings.Split(out, "\n")
	if len(lines) != 8 {
		t.Errorf("expected 8 lines, got %d", len(lines))
	}
	if strings.HasSuffix(out, "\n") {
		t.Error("trailing newline should be stripped")
	}
	if !strings.Contains(out, "hi") || !strings.Contains(out, "hello") {
		t.Error("expected message text in output")
	}
}

func TestMessageList_Scroll(t *testing.T) {
	ml := NewMessageList()
	msgs := []models.Message{
		{UserID: "a", Text: "m1"},
		{UserID: "b", Text: "m2"},
		{UserID: "c", Text: "m3"},
		{UserID: "d", Text: "m4"},
	}
	for _, m := range msgs {
		ml.AddMessage(m)
	}

	out0 := ml.Render(50, 4)
	ml.ScrollUp(2)
	out1 := ml.Render(50, 4)
	if out0 == out1 {
		t.Error("scroll should change output")
	}

	ml.ScrollDown(2)
	out2 := ml.Render(50, 4)
	if out0 != out2 {
		t.Error("scroll back should match original")
	}
}

func TestRoomList_Render(t *testing.T) {
	rl := NewRoomList([]string{"Room A", "Room B"})
	out := rl.Render(50)
	if !strings.Contains(out, "Room A") || !strings.Contains(out, "Room B") {
		t.Error("expected room names in output")
	}
	lines := strings.Split(out, "\n")
	if len(lines) != 3 {
		t.Errorf("expected 3 lines (2 rooms + trailing), got %d", len(lines))
	}
}

func TestRoomList_Navigation(t *testing.T) {
	rl := NewRoomList([]string{"Room A", "Room B", "Room C"})
	if rl.SelectedRoom != 0 {
		t.Error("should start at index 0")
	}
	rl.OnKeyPress("down")
	if rl.SelectedRoom != 1 {
		t.Error("down should increment")
	}
	rl.OnKeyPress("down")
	if rl.SelectedRoom != 2 {
		t.Error("down should stop at last")
	}
	rl.OnKeyPress("down")
	if rl.SelectedRoom != 2 {
		t.Error("down should not go past end")
	}
	rl.OnKeyPress("up")
	if rl.SelectedRoom != 1 {
		t.Error("up should decrement")
	}
	rl.OnKeyPress("up")
	if rl.SelectedRoom != 0 {
		t.Error("up should stop at first")
	}
	rl.OnKeyPress("up")
	if rl.SelectedRoom != 0 {
		t.Error("up should not go before start")
	}
	redirect, _ := rl.OnKeyPress("enter")
	if redirect != RedirectSelect {
		t.Error("enter should return RedirectSelect")
	}
}

func TestUserList_AddRemoveCount(t *testing.T) {
	ul := NewUserList()
	if ul.Count() != 0 {
		t.Error("should start empty")
	}
	ul.AddUser("user1")
	if ul.Count() != 1 {
		t.Error("add should increment")
	}
	ul.AddUser("user1")
	if ul.Count() != 1 {
		t.Error("duplicate add should be ignored")
	}
	ul.AddUser("user2")
	if ul.Count() != 2 {
		t.Error("second user should increment")
	}
	ul.RemoveUser("user1")
	if ul.Count() != 1 {
		t.Error("remove should decrement")
	}
	ul.RemoveUser("nonexistent")
	if ul.Count() != 1 {
		t.Error("remove nonexistent should not change count")
	}
}

func TestTitleBar_Render(t *testing.T) {
	tb := NewTitleBar("TEST")
	out := tb.Render(30)
	if !strings.Contains(out, "TEST") {
		t.Error("expected title text")
	}
	outLines := strings.Split(out, "\n")
	if len(outLines) != 1 {
		t.Errorf("title should be 1 line, got %d", len(outLines))
	}
}

func TestFooter_Render(t *testing.T) {
	f := NewFooter("help text")
	out := f.Render(50)
	if !strings.Contains(out, "help text") {
		t.Error("expected footer text")
	}
}

func TestUserPanel_RenderFull(t *testing.T) {
	up := NewUserPanel()
	out := up.RenderFull(10)
	lines := strings.Split(out, "\n")
	if len(lines) != 10 {
		t.Errorf("expected 10 lines, got %d", len(lines))
	}
	if !strings.Contains(out, "ONLINE") {
		t.Error("expected ONLINE header")
	}
}

func TestMessageList_Render_ExactFit(t *testing.T) {
	ml := NewMessageList()
	ml.AddMessage(models.Message{UserID: "x", Text: "a"})
	out := ml.Render(50, 2)
	lines := strings.Split(out, "\n")
	if len(lines) != 2 {
		t.Errorf("expected 2 lines, got %d", len(lines))
	}
}
