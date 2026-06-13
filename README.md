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

### Custom rooms

```bash
go run ./tui --rooms "Lobby,Go,Random"
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
- `tui/ui/` — screen implementations (room list, chat room)
- `tui/components/` — UI widgets and shared lipgloss palette
- `.github/workflows/go.yml` — CI pipeline
