# Terminal Chat Application

A terminal-based chat application with a modern TUI interface, supporting multiple chat rooms and real-time messaging via Redis.

## Features

- Terminal-based user interface with Bubble Tea
- Multiple chat rooms
- Real-time messaging using Redis pub/sub
- SSH server for remote access
- User identification and message coloring

## Prerequisites

- Go 1.21 or higher
- Redis server
- Git (for cloning the repository)

## Installation

### 1. Clone the Repository

```bash
git clone <repository-url>
cd terminal_chat
```

### 2. Install Dependencies

```bash
# Install server dependencies
cd server
go mod download

# Install TUI dependencies
cd ../tui
go mod download
```

### 3. Build the Application

```bash
# Build the TUI client
cd tui
go build -o chat

# Build the server
cd ../server
go build -o chat-server
```

## Server Setup

### 1. Redis with Docker

You can easily run Redis using Docker:

```bash
# Pull the latest Redis image
docker pull redis

# Run Redis container
docker run --name redis-chat -p 6379:6379 -d redis

# To persist data, use a volume:
# docker run --name redis-chat -p 6379:6379 -v redis_data:/data -d redis

# Verify the container is running
docker ps

# View logs
docker logs redis-chat

# Stop the container when needed
docker stop redis-chat

# Start it again later
docker start redis-chat
```

Or use Docker Compose (create a `docker-compose.yml` file):

```yaml
version: '3.8'
services:
  redis:
    image: redis:alpine
    container_name: redis-chat
    ports:
      - "6379:6379"
    volumes:
      - redis_data:/data
    restart: unless-stopped

volumes:
  redis_data:
```

Then run:
```bash
docker-compose up -d
```

### 2. Server Configuration

The server can be configured using environment variables:

| Variable | Default | Description |
|----------|---------|-------------|
| `HOST` | `0.0.0.0` | Host to bind the server to |
| `PORT` | `2222` | Port to listen on |
| `REDIS_ADDR` | `localhost:6379` | Redis server address |
| `REDIS_PASSWORD` | `` | Redis password (if required) |
| `REDIS_DB` | `0` | Redis database number |

### 3. Running the Server

```bash
# Navigate to server directory
cd server

# Run the server (with default settings)
./chat-server

# Or with custom settings
HOST=0.0.0.0 PORT=2222 ./chat-server
```

### 4. Running as a Service (Systemd)

Create a systemd service file at `/etc/systemd/system/chat-server.service`:

```ini
[Unit]
Description=Terminal Chat Server
After=network.target

[Service]
Type=simple
User=your_username
WorkingDirectory=/path/to/terminal_chat/server
Environment="PATH=/usr/local/sbin:/usr/local/bin:/usr/sbin:/usr/bin:/sbin:/bin"
ExecStart=/path/to/terminal_chat/server/chat-server
Restart=always

[Install]
WantedBy=multi-user.target
```

Then enable and start the service:

```bash
sudo systemctl daemon-reload
sudo systemctl enable chat-server
sudo systemctl start chat-server
```

## Client Usage

### Connect to the Server

```bash
ssh username@server-address -p 2222
```

### Keyboard Controls

- `↑/↓`: Navigate room list
- `Enter`: Select room/send message
- `Esc`: Go back to room list
- `Ctrl+C` or `Ctrl+D`: Exit the application

## Security Considerations

1. **Authentication**: The default setup has no authentication. For production use:
   - Set up SSH key-based authentication
   - Or implement password authentication in the server

2. **Firewall**: Ensure your firewall allows the SSH port (default: 2222)

3. **TLS**: Consider using a reverse proxy like Nginx with Let's Encrypt for encrypted connections

## Development

### Project Structure

- `/tui`: Terminal user interface code
- `/server`: SSH server implementation
- `/tui/messageStore`: Redis client and message handling
- `/tui/models`: Data models and state management
- `/tui/ui`: User interface components

### Adding Features

1. Clone the repository
2. Create a new branch
3. Make your changes
4. Test thoroughly
5. Submit a pull request

## License

MIT License

## Contributing

Contributions are welcome! Please feel free to submit a Pull Request.
