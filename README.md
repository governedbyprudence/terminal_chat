# Terminal Chat

A terminal-based chat application with a modern TUI interface supporting multiple chat rooms.

## Features

- Terminal user interface with [Bubble Tea](https://github.com/charmbracelet/bubbletea)
- Multiple chat rooms
- Real-time messaging via embedded TCP server (no external dependencies)
- User identification and message coloring

## Usage

### Standalone (single user, no networking)

```bash
go run ./tui
```

### Multi-user over TCP

Terminal 1 (server + TUI):
```bash
go run ./tui --server :9000
```

Terminals 2, 3, ... (clients):
```bash
go run ./tui --connect localhost:9000
```

## Controls

- `↑/↓`: Navigate room list
- `Enter`: Select room / send message
- `Esc`: Go back to room list
- `q`: Quit

## Project Structure

- `tui/main.go` — entry point, program loop, message subscription
- `tui/messageStore/` — TCP pub/sub broker (also handles local/in-memory mode)
- `tui/models/` — data models and state management
- `tui/ui/` — UI components (room list, chat room)
- `tui/styles/` — lipgloss styling
- `tui/utils/` — user ID utilities
