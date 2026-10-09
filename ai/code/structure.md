# Code Structure

## Overview

Terminal Chat is a Go 1.24 TUI chat application using [Bubble Tea](https://github.com/charmbracelet/bubbletea) v1.3.6, [Lip Gloss](https://github.com/charmbracelet/lipgloss) v1.1.0, and [Bubbles](https://github.com/charmbracelet/bubbles) v0.21.0. Supports single-user (in-memory pub/sub) and multi-user (TCP broker) modes across multiple chat rooms.

**Source:** ~2,000 Go LOC across 19 `.go` files (+3 `_test.go`) in 5 packages (`main`, `ui`, `components`, `messageStore`, `models`).

---

## 1. Package Analysis

### Dependency Map

```
main ──→ ui ──→ components ──→ models
  │          │
  │          ├──→ messageStore
  │          └──→ models
  │
  ├──→ components
  ├──→ messageStore
  └──→ models
```

### Package State

| Package | Content | Notes |
|---------|---------|-------|
| `models` | `Message`, `PresenceEvent`, `AppScreen` interface, `CurrentUserID()` | Types + user ID generation. |
| `messageStore` | `Broker` interface + unexported `broker` (in-memory pub/sub, TCP server/client) | Interface hides transport. Channel names are typed constants. TCP server/client support for multi-user mode. 6 tests. |
| `components` | All UI widgets (`InputBar`, `MessageList`, `RoomList`, `RoomTabs`, `UserList`, `TitleBar`, `Footer`, `UserPanel`) + `Component` interface + shared palette | Widgets own inline styles; shared menu styles in `palette.go`. `MessageList` exposes `ScrollTrack()` for the scroll indicator column. 11 tests. |
| `ui` | `ChatRoomScreen` + `layout.go` (shared outer frame + `sepLine`) | Single-page app. `ChatRoomScreen` receives `Broker` via dependency injection. Body height clamped to prevent overflow. 10 tests. |
| `main` | Entry point, Bubble Tea `Model`, `qq` quit handler | Subscriptions created once in `newModel()`, long-lived. `awaitingQuit bool` + `quitResetMsg` for double-tap quit (800ms window). |

---

## 2. Architecture

```
                     main.go
                      Model
                     /      \
            messageStore    ui/ (screens)
                Broker     └──ChatRoomScreen (single-page)
             ─────────                  |
               broker          components/
                             ┌──────────────────┐
                             │   Palette        │── shared colors + styles
                             ├──────────────────┤
                             │  Component       │── interface
                             ├──────────────────┤
                             │  TitleBar        │
                             │  Footer          │
                             │  InputBar        │
                             │  MessageList     │
                             │  RoomList        │
                             │  RoomTabs        │── tab bar (neo-brutalist)
                             │  UserList        │
                             │  UserPanel       │── composite wrapping UserList
                             └──────────────────┘
```

### Key Design Decisions

- **Component styles**: Unique widget styles defined inline per file; shared menu styles (`MenuItem`, `SelectedMenuItem`) consolidated in `palette.go` alongside the 5-color palette.
- **`RoomTabs` component**: Horizontal tab bar rendered with `lipgloss.JoinHorizontal`. Active tab gets `Accent` background + `DeepBg` text + bold. Inactive tabs get `Text` foreground. Room names truncate with `…` if they overflow the available width. Supports `SelectNext()`, `SelectPrev()`, `SelectIndex(i)`.
- **`ChatRoomScreen` owns all rooms**: The screen holds a `Rooms []string` slice and per-room state maps: `messageLists map[string]*MessageList` (each room keeps its own scroll position) and `usersPerRoom map[string]map[string]bool` (per-room user presence sets). Tab switching (`←`/`→`/`tab`/`shift+tab`) publishes presence leave/join events.
- **`UserPanel` composite**: Encapsulates the right sidebar — "ONLINE" header, user list, and vertical fill. Repopulated on tab switch via `syncUserPanel()`.
- **Single-page app**: `ChatRoomScreen` is the only screen. The tab bar IS the navigation. ESC is a no-op (no back reference).
- **Fixed-height layout**: `ChatRoomScreen.Render(width, height)` partitions the terminal into `contentH = height - 2` lines. Layout: tabs (1) + separator (1) + combined left/right panel (`belowSepH = contentH - 2`). The left panel (contentW - RightPanelWidth - PanelGap) contains messages + input bar (3) + footer (1). A 2-char scroll track column sits between the left panel and the right user panel, inserted via `JoinHorizontal`. Combined lines are padded per-line to `contentW` (bypasses lipgloss wrapper bug with multi-line ANSI content). A body-height clamp prevents overflow on very small terminals.
- **Scroll track indicator**: When messages overflow the viewport, a neo-brutalist scroll indicator appears as a 2-char column to the left of the user panel divider. Uses `░░` (light shade) for the track bed and `██` (full block) for the thumb. Exposed via `MessageList.ScrollTrack(maxHeight) → []string` — proportional thumb position calculated from `scrollOffset / maxOffset`. Built in `ChatRoomScreen.Render` by calling `ScrollTrack(msgHeight)`, padding to `belowSepH` height, and joining between left content and user panel. `msgHeight` is cached on the screen struct for use in `Update`.
- **Input bar border alignment**: The input bar's top/bottom borders (`inputWrapper` with `Accent`-colored `─`) are extended 2 chars to the right using `extStyle.Render("──")` in the scroll track column, so the border aligns with the scrollbar column's right edge.
- **Scroll key bindings**: `↑`/`↓` scrolls one message at a time. `PgUp`/`PgDn` scrolls by a fixed 3 messages. `ChatRoomScreen.Update` intercepts these before the default handler (which forwards to `InputBar`).
- **`qq` double-tap quit**: Single `q` is intercepted and sets `awaitingQuit = true` with an 800ms reset timer via `tea.Tick`. A second `q` before the timer fires sends `tea.Quit`. Any other key or timer expiry resets the flag. Prevents accidental quits during typing.
- **`ui/layout.go`**: Shared `outerFrame` style (double border, accent color) and `sepLine(width)` helper. `outerFrameVSize()` = 2, `outerFrameHSize()` = 6.
- **Broker**: Interface (`Subscribe`, `Publish`, `Connect`, `Disconnect`, `StartServer`, `StopServer`) with unexported concrete type. Passed from `main()` → `ChatRoomScreen`. Two pub/sub channels: `"chat"` for messages, `"presence"` for join/leave events. Wire format is JSON.

---

## 3. Test Infrastructure

27 tests across 3 packages, run with `go test ./tui/...`:

| Package | File | Tests | Scope |
|---------|------|-------|-------|
| `components` | `components_test.go` | 11 | `MessageList` render (empty, non-empty, scroll, exact-fit), `RoomList`, `UserList`, `TitleBar`, `Footer`, `UserPanel` rendering |
| `messageStore` | `messageStore_test.go` | 6 | Publish/subscribe delivery, multiple subscribers, channel isolation, close cleanup, TCP server → client integration |
| `ui` | `ui_test.go` | 10 | `ChatRoomScreen` line count, small-terminal clamp, message routing, pgup/pgdown scroll, enter publish, `RoomListScreen` render + enter |

## 4. Design Reference

See [ai/design/brutalist-layout.md](../design/brutalist-layout.md) for the neo-brutalist design system spec, color palette, typography rules, and the full room tab integration rationale.
