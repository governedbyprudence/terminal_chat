# Brutalist Layout: bubbletea + lipgloss Lessons

## Cursor Overwrites Last Frame Character

When using `tea.WithAltScreen()`, the `View()` output **must end with `"\n"`**. Without it, the terminal cursor position overlaps the last character of the last rendered line — typically the frame's bottom-right corner (`╝`).

```go
// WRONG — corner missing:
return styles.OuterFrame.Render(body)

// RIGHT — corner visible:
return styles.OuterFrame.Render(body) + "\n"
```

## Bubbles TextInput Width Bug

The `textinput.Model` always renders **Width + 1** visual characters. The `View()` method has a `padding++` bug that adds an extra space beyond the configured `Width`. This holds regardless of whether the cursor is on a character or at the end.

```go
// If you want the textinput to visually occupy W characters:
ti.Width = W - 1

// In a chat layout where contentW = terminalWidth - frameOverhead:
ti.Width = contentW - 7
```

Breakdown: `prompt(2) + textinput(W+1) + wrapperPadding(4) = W + 7 = contentW`.

## TextStyle Padding Is Silently Ignored

The textinput internally calls `m.TextStyle.Inline(true).Render()`. The `Inline(true)` flag causes lipgloss to skip padding, borders, and margins. So `Padding(0, 1)` on `TextStyle` has zero effect:

```go
TextInputStyle := lipgloss.NewStyle().
    Foreground(Text).
    Padding(0, 1)  // ← IGNORED by textinput
```

Remove padding from `TextStyle`; add visual spacing via a wrapper style instead.

## `Width()` Is Minimum, Not Exact

`style.Width(N)` sets a **minimum** width. If content is already wider, it stays wider. You cannot use it to truncate or to force an exact output width.

## `Width()` With Alignment Breaks Manual Padding

A style with `Width(N)` and `Align(Right)` pads short lines to N using the alignment direction. If you've already manually right-padded a line (spaces prepended with `strings.Repeat`), applying this style re-pads the line with spaces on the left again, destroying your manual alignment.

**Fix:** use a plain wrapper with no alignment set (defaults to left, which is a no-op for lines already at minimum width):

```go
wrapper := lipgloss.NewStyle().Width(contentW)
wrapped := wrapper.Render(body)
```

## Use `strings.Repeat` for Reliable Right-Alignment

Manual padding is immune to style composition issues:

```go
rightPad := func(s string) string {
    w := lipgloss.Width(s)
    if w >= contentW {
        return s
    }
    return strings.Repeat(" ", contentW-w) + s
}
```

Use `lipgloss.Width()` (not `len()`) to measure visual width — it strips ANSI escape codes correctly.

## Frame Width Arithmetic

For a frame with `Padding(0, 2)` and a 1-char border on each side:

```
OuterFrameHSize = 1(border-left) + 1(border-right) + 2(padding-left) + 2(padding-right) = 6
contentW = terminalWidth - OuterFrameHSize
```

The widest content line determines frame width. If all content lines ≤ contentW, the frame is exactly `terminalWidth`. Any line wider than contentW causes overflow (missing right border).

## Debug Logging in TUI

Use `log.Println()` for debug output — it goes to stderr, visible in the terminal that started the TUI app (not on the TUI screen):

```go
log.Printf("frame width=%d, terminal width=%d", frameW, width)
```

Only the model's `View()` string appears on screen; everything else goes to stderr or the log file.

# Component Architecture

## Motivation

The original code had two screens (`RoomListScreen`, `ChatRoomScreen`) with duplicated framing, inline rendering closures (the `rightPad` helper was redefined on every `Render` call), and hard-to-follow navigation with type assertions all over `main.go`.

We extracted patterns into reusable, composable units.

## The `Component` Interface

Every reusable piece implements this contract (`tui/components/component.go`):

```go
type Component interface {
    Render(width int) string
    OnKeyPress(key string) (redirect int, cmd tea.Cmd)
    Focus() tea.Cmd
    Blur()
    Focused() bool
}
```

- `Render(width)` — draw the component at a given content width.
- `OnKeyPress` — handle a key; return `redirect=1` to signal navigation.
- `Focus/Blur/Focused` — focus management for interactive components.

Non-interactive components (e.g. `TitleBar`, `Footer`) implement the interface with no-ops — this avoids type-assertion checks at call sites while keeping the interface uniform.

## Presentational Components

These are stateless or have simple public state. They ignore `width` (mark with `_`) since their content determines their size.

| Component | File | Purpose |
|---|---|---|
| `TitleBar` | `titlebar.go` | Styled header via `styles.TitleBar` |
| `Footer` | `footer.go` | Dimmed navigation hint via `styles.FooterText` |
| `MessageList` | `messagelist.go` | Renders sent/received messages with right-aligned incoming text. Owns a `[]states.Message` slice. |

## Interactive Components

These manage internal state and respond to keys.

| Component | File | Behavior |
|---|---|---|
| `InputBar` | `inputbar.go` | Wraps `textinput.Model`. Exposes `Value()`, `SetValue()`, `Update(tea.Msg)` for text processing. Sets `Width = contentW - 7` per the bubbles bug workaround. |
| `RoomList` | `roomlist.go` | Manages `SelectedRoom` index. `OnKeyPress("down"/"up")` adjusts selection; `"enter"` returns `redirect=1`. |

## The `BorderFrame` Decorator

`BorderFrame` (`borderframe.go`) wraps any `Component` with the outer frame and trailing newline:

```go
type BorderFrame struct {
    Inner Component
}

func (f *BorderFrame) Render(width int) string {
    contentW := width - styles.OuterFrameHSize  // 6 = border(2) + padding(4)
    body := f.Inner.Render(contentW)
    wrapper := lipgloss.NewStyle().Width(contentW)
    return styles.OuterFrame.Render(wrapper.Render(body)) + "\n"
}
```

It delegates `OnKeyPress`, `Focus`, `Blur`, `Focused` to the inner component — a textbook decorator. This ensures every framed view always gets the trailing newline (fixing the cursor-overlap bug).

## Screen Composition

Screens are NOT components. They implement `AppScreen` (`tui/models/states/states.go`):

```go
type AppScreen interface {
    Render(width, height int) string
    Update(msg tea.Msg) (AppScreen, tea.Cmd)
}
```

Screens own component instances and wire them together:

```go
type ChatRoomScreen struct {
    RoomName    string
    MessageList *components.MessageList
    InputBar    *components.InputBar
    TitleBar    *components.TitleBar
    Footer      *components.Footer
    backRef     states.AppScreen  // navigation target on "esc"
}
```

```go
type RoomListScreen struct {
    TitleBar *components.TitleBar
    RoomList *components.RoomList
    Footer   *components.Footer
}
```

Both screens handle their own framing in `Render` (or could use `BorderFrame` for delegation) — the key requirement is every `View()` ending with `"\n"`.

## Navigation Without Type Assertions

The original `main.go` used `m.Screen.(ui.ChatRoomScreen)` to access the textinput and route messages. The refactored code uses a `backRef` field:

- `ChatRoomScreen` stores a reference to the room list screen when constructed.
- On `"esc"`, its `Update` returns `s.backRef` directly.
- `RoomListScreen` creates a fresh `ChatRoomScreen` on `"enter"`, passing itself as the back reference.

```go
// RoomListScreen creates ChatRoomScreen with self as backRef
case "enter":
    chat := NewChatRoomScreen(s, roomName)
    return chat, textinput.Blink

// ChatRoomScreen returns to RoomListScreen
case "esc":
    return s.backRef, nil
```

The main model is now three cases — `WindowSizeMsg`, `KeyMsg` (delegates to screen), `states.Message` (delegates to screen). No type assertions.

## Lessons

1. **Every `View()` must end with `"\n"`** — both screens need it, not just one. A missing newline causes the cursor to eat the last frame character after switching screens.
2. **Components as decorators** — `BorderFrame` wraps any component with framing. This ensures consistent application of width arithmetic and the trailing-newline rule.
3. **Uniform interface with no-op stubs** — non-interactive components implement `Focus/Blur/Focused` as no-ops rather than splitting the interface. This avoids type assertions at call sites.
4. **`backRef` over screen managers** — for a two-screen app, a simple back-reference in the constructor is cleaner than a routing layer. The screens own their transitions.
5. **Width is a minimum** — components that don't care about width use `_` for the parameter. Components that do (like `MessageList` for right-padding) use it actively.
6. **State lives in components** — `RoomListState` and `ChatState` were removed; the `RoomList` and `InputBar` components manage their own state internally.
