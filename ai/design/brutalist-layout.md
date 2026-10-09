# Brutalist Layout: bubbletea + lipgloss Lessons

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

**Fix:** pad each line individually so alignment is not re-applied:

```go
wrapper := lipgloss.NewStyle().Width(contentW)
wrapped := wrapper.Render(body)
```

**Caveat:** `style.Width(N).Render(multiLine)` has a lipgloss bug where rendering multi-line content containing ANSI-coded `JoinHorizontal` output drops content from the first line. Workaround: pad each line manually instead of using a wrapper for multi-line bodies.

## Use `strings.Repeat` for Reliable Right-Alignment and Minimum Width

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

For enforcing a minimum line width on multi-line ANSI content (avoids the wrapper line-dropping bug):

```go
lines := strings.Split(body, "\n")
for i, line := range lines {
    if w := lipgloss.Width(line); w < minWidth {
        lines[i] = line + strings.Repeat(" ", minWidth-w)
    }
}
body = strings.Join(lines, "\n")
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
    Rooms         []string
    ActiveRoomIdx int
    InputBar      *components.InputBar
    Footer        *components.Footer
    UserPanel     *components.UserPanel
    RoomTabs      *components.RoomTabs
    messageLists  map[string]*components.MessageList
    usersPerRoom  map[string]map[string]bool
    broker        messageStore.Broker
    backRef       models.AppScreen
}
```

```go
type RoomListScreen struct {
    TitleBar *components.TitleBar
    RoomList *components.RoomList
    Footer   *components.Footer
    broker   messageStore.Broker
}
```

Screens handle their own framing in `Render`.

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

1. **Uniform interface with no-op stubs** — non-interactive components implement `Focus/Blur/Focused` as no-ops rather than splitting the interface. This avoids type assertions at call sites.
2. **`backRef` over screen managers** — for a two-screen app, a simple back-reference in the constructor is cleaner than a routing layer. The screens own their transitions.
3. **Width is a minimum** — components that don't care about width use `_` for the parameter. Components that do (like `MessageList` for right-padding) use it actively.
4. **State lives in components** — `RoomListState` and `ChatState` were removed; the `RoomList` and `InputBar` components manage their own state internally.

---

# Neo-Brutalist Design System

## Design Philosophy

Neo-brutalism in terminal UI means: maximal contrast, heavy borders, raw typography, and zero decoration. Every pixel earns its place. The UI communicates hierarchy through weight (thickness, boldness, color intensity) rather than chrome (gradients, shadows, icons).

## Color Palette

| Token | Hex | Usage |
|---|---|---|
| `DeepBg` | `#0A0A0A` | Terminal background — near-black, absorbs light |
| `Surface` | `#161616` | Component surfaces, list backgrounds |
| `Text` | `#F0F0F0` | Primary text — bright, high contrast |
| `DimText` | `#585858` | Secondary text, timestamps, hints |
| `Accent` | `#FF5E00` | Vivid orange — active tab, highlights, selected item, border color |
| `Danger` | `#E03030` | Leave/quit actions, errors (future use) |

The palette is intentionally constrained to 5 colors (+1 secondary). This forces deliberate use of weight, spacing, and layout rather than color-coding to convey meaning.

## Typography

- Monospace only (terminal-native).
- Bold weight for: active tabs, room names, sent message labels, title bars.
- Normal weight for: received messages, user lists, metadata.
- Italic for: timestamps, empty states ("no messages yet", "nobody here").
- Dim text for: footer hints, secondary metadata.

## Layout Principles

1. **Full-height framing** — every screen fills the terminal exactly. Padding/clamping ensures no dead space.
2. **Heavy borders** — `DoubleBorder` (unicode box-drawing `╔═╗`) for the outer frame in accent orange. `NormalBorder` for internal dividers and separators.
3. **Tab-based navigation** — rooms are tabs at the top of the chat screen. No separate room-list screen. The tab bar IS the navigation.
4. **Right-align inbound messages** — visual cue that "they" are on the other side. Outbound messages are left-aligned.
5. **User panel on the right** — bordered left edge, fixed 22-char width. Shows who's in the current room.
6. **Input bar at the bottom** — bordered top edge, always visible, always focused.

## Tab Bar Specification

Room tabs sit directly below the outer frame's top border, spanning the full `contentW`. No separate title bar — the active tab itself identifies the room:

```
╔═══════════════════════════════════════════╗
║  ╭ Room 1 ╮  ╭ Room 2 ╮  ╭ Room 3 ╮     ║  ← tabs row (╭ ╮ rounded corners)
║  ───────────────────────────────────────  ║  ← sepLine (─ repeated)
║  messages...                              ║
║  messages...                              ║
║  ───────────────────────────────────────  ║  ← input bar top border
║  > hello                                  ║
╚═══════════════════════════════════════════╝
```

**Active tab** — Accent background (`#FF5E00`), `DeepBg` foreground, bold, with rounded-corner characters:
```
╭ Room 2 ╮
```

**Inactive tab** — Surface background (`#161616`), `Text` foreground, normal weight:
```
╭ Room 1 ╮
```

**Keyboard navigation:**
- `←` / `→` — switch to adjacent room
- `tab` — switch to next room
- `shift+tab` — switch to previous room

## Rationale for Single-Screen Design

Eliminating the RoomListScreen:

1. **Reduces cognitive load** — users don't navigate away from their conversation to switch rooms.
2. **Matches modern chat UX** — Slack, Discord, Telegram all use channel tabs/rails.
3. **Faster room switching** — one keystroke instead of `esc` → navigate → `enter`.
4. **Presence continuity** — leaving a room to pick another no longer sends a "leave" event.

The tradeoff: the initial experience is slightly busier (all rooms shown at once). But for a terminal chat app, power users will prefer the speed.

## Message Layout

```
you              ← accent "you" label, bold
hey team!        ← sent message, bold
                                   received message →  ← right-padded to full width
                                              bob →  ← dim sender label, right-aligned
```

- **Sent messages** (current user): left-aligned, bold, with `you` label above.
- **Received messages** (others): right-aligned (using `rightPad`), normal weight, with sender short-ID below right-aligned.
- Each message occupies 2 lines (label line + text line).

## Tab Bar Implementation

The tab bar is rendered as a row using `lipgloss.JoinHorizontal(lipgloss.Top, ...)`.

Each tab is a styled string:

```go
func renderTab(name string, active bool, width int) string {
    if active {
        return activeTab.Render("╭ " + name + " ╮")
    }
    return inactiveTab.Render("╭ " + name + " ╮")
}
```

Tab widths are calculated to fit the available content width. If tabs overflow, they are truncated with `…`.


