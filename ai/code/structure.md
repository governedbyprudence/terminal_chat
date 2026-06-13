# Code Structure & Improvement Plan

## Overview

Terminal Chat is a Go 1.24 TUI chat application using [Bubble Tea](https://github.com/charmbracelet/bubbletea) v1.3.6, [Lip Gloss](https://github.com/charmbracelet/lipgloss) v1.1.0, and [Bubbles](https://github.com/charmbracelet/bubbles) v0.21.0. It supports single-user (in-memory pub/sub) and multi-user (TCP broker) modes across multiple chat rooms.

**Source:** ~1,400 Go LOC across 15 `.go` files in 5 packages (`main`, `ui`, `components`, `messageStore`, `models`).

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
| `models` | `Message`, `PresenceEvent`, `AppScreen` interface, `CurrentUserID()` | Types + user ID generation. Pure types + one helper. |
| `messageStore` | `Broker` interface + unexported `broker` implementation (in-memory pub/sub, TCP server/client) | Exports clean interface, hides transport details. Channel names are typed constants. |
| `components` | All UI widgets (`InputBar`, `MessageList`, `RoomList`, `UserList`, `TitleBar`, `Footer`) + `UserPanel` composite + `Component` interface + shared palette | Each component file owns its lipgloss styles inline. Palette colors in `palette.go`. `UserPanel` encapsulates the right sidebar (header, user list, visual gap, vertical stretch). |
| `ui` | Screen implementations (`RoomListScreen`, `ChatRoomScreen`) + `layout.go` (shared outer frame style) | Screens receive `Broker` via dependency injection. No global state calls. |
| `main` | Entry point, Bubble Tea `Model`, subscription lifecycle | Subscriptions created once in `newModel()`, long-lived. |

---

## 2. Architecture

```
                    main.go
                     Model
                    /      \
           messageStore    ui/ (screens)
               Broker     ┌──RoomListScreen
            ─────────     └──ChatRoomScreen
            broker                |
                         components/
                         ┌──────────────┐
                         │   Palette    │── shared 5-color palette
                         ├──────────────┤
                         │  Component   │── interface (Render, OnKeyPress, Update, Focus/Blur/Focused)
                         ├──────────────┤
                         │  TitleBar    │
                         │  Footer      │
                         │  InputBar    │
                         │  MessageList │
                         │  RoomList    │
                         │  UserList    │
                         │  UserPanel   │── composite wrapping UserList
                         └──────────────┘
```

### Key Design Decisions

- **Self-contained components**: Each component owns its lipgloss styles inline. The 5-color palette in `components/palette.go` is the only shared style dependency. No separate `styles/` package.
- **`UserPanel` composite**: Encapsulates the entire right sidebar — header ("ONLINE"), user list, visual gap separator, and vertical padding to match left panel height. Screens no longer do inline layout arithmetic.
- **`Component` interface**: Consistent API across all widgets (`Render`, `OnKeyPress`, `Update`, `Focus`/`Blur`/`Focused`). `InputBar.Update` is now part of the interface rather than a standalone method.
- **`ui/layout.go`**: Shared `outerFrame` style (rounded border with accent) used by both screens, avoiding duplication. `outerFrameVSize()` returns 2 (top + bottom border lines), `outerFrameHSize()` returns 4 (left + right border columns).
- **Fixed-height layout mechanics**: `ChatRoomScreen.Render(width, height)` partitions the terminal into a body of exactly `contentH = height - 2` lines. Body content is split: `msgHeight = contentH - 5` lines for the message list (which uses bottom-aligned padding internally), 1 title line, 1 blank separator, 3 input bar lines, and 1 footer line. The outer frame wraps the body to exactly `height` lines — no trailing newlines, no overflow. Past off-by-one bugs from trailing `\n` in rendered strings and word-wrapping of concatenated footer content have been fixed via `strings.TrimRight` and proper line separation.
- **Broker interface** in `messageStore` (with methods `Subscribe`, `Publish`, `Connect`, `StartServer`, `StopServer`) — concrete `broker` struct is unexported, constructed via `New()`.
- **Dependency injection**: Broker is created in `main()`, passed through `newModel()` → `RoomListScreen` → `ChatRoomScreen`.
- **Long-lived subscriptions**: Created once in `newModel()`, stored on `Model` as `*Subscription` pointers. Value-receiver copies share the same underlying subscriptions.
- **Channel name constants**: `messageStore.ChannelChat = "chat"`, `messageStore.ChannelPresence = "presence"`.
- **Self-documenting wire protocol**: `wireMsg` uses `json:"channel"` and `json:"payload"`.

---

## 3. Remaining Issues

### 3.1 No Test Infrastructure

Zero `*_test.go` files exist. The entire application has:
- No unit tests for the broker
- No unit tests for components
- No integration tests for the TCP protocol
- No TUI model tests

### 3.2 Over-Exported Names

Types like `TitleBar`, `Footer`, `RoomList`, `MessageList`, `UserList`, `InputBar`, `UserPanel` are exported but only used within their package or by `main.go`. Only `AppScreen` is a genuinely public interface.

### 3.3 Inconsistent API Surface

- `InputBar.Render()` calls `Focus()` on every render, unlike other components.
- `ChatRoomScreen` calls `s.InputBar.Update(msg)` directly rather than dispatching through the `Component` interface.

### 3.4 `Component` Interface Has No Polymorphic Consumers

All components implement `Component`, but nothing depends on the interface for dispatch — screens access children by concrete type. The actual polymorphism runs through `AppScreen`.

### 3.5 `readLoop` Cannot Be Externally Stopped

No context, cancel channel, or close signal. The goroutine runs until the connection errors or the process exits.

### 3.6 Error Handling Edge Cases

| Location | Issue | Severity |
|----------|-------|----------|
| `messageStore.go:110` | `reader.ReadString` error terminates read loop silently | Major |
| `messageStore.go:115` | `json.Unmarshal` error silently `continue`s | Minor |
| `messageStore.go:165` | `w.WriteString` error silently `continue`s | Minor |

---

## 4. Improvement Plan

### Phase 1: Add Test Infrastructure

| Task | Scope | Status |
|------|-------|--------|
| 1.1 Write broker tests | `messageStore/messageStore_test.go` — `Publish`/`Subscribe` delivery, channel isolation, `Close()` cleanup, buffered non-blocking. | ✅ Done |
| 1.2 Write component tests | `components/components_test.go` — `MessageList.Render` (empty, non-empty, scroll, exact-fit), `RoomList` navigation, `UserList` add/remove, `TitleBar`/`Footer`/`UserPanel` rendering. 10 tests total. | ✅ Done |
| 1.3 Write screen functional tests | `ui/ui_test.go` — ChatRoomScreen line count with and without messages, small-terminal clamp, no trailing NL, message routing, esc navigation, pgup/pgdown scroll, enter publish, RoomListScreen render + navigation. 10 tests total. | ✅ Done |
| 1.4 Write integration test | TCP server → multiple clients → message delivery. | ❌ Pending |

### Phase 2: Code Quality & Maintenance

| Task | Scope |
|------|-------|
| 2.1 Add Go CI workflow | `.github/workflows/go.yml` — `go test ./...` on push/PR. |
| 2.2 Make room names configurable | `main.go` — Accept via CLI flag or env var with sensible defaults. |
| 2.3 Integrate or remove `Component` interface | Either use it for polymorphic rendering in screens (e.g. `[]Component` for iterating Render/Update) or remove the unused interface methods. |
| 2.4 Address `readLoop` stop | Add context or close channel to `readLoop`. |
| 2.5 Add project-local `.gitignore` for `debug.log` | Already exists in `tui/.gitignore` but `debug.log` was committed before. |
