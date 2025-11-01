# Go Real-Time Chat Application

A high-performance, real-time chat application built with Go and WebSockets. This chat application supports multiple rooms, user presence tracking, and provides a modern web interface.

## Features

- 🔌 **Real-Time Messaging**: Instant message delivery using WebSocket connections
- 🏠 **Multiple Rooms**: Create and join different chat rooms
- 👥 **User Presence**: See who joins and leaves in real-time
- 📊 **Room Statistics**: Track active users per room
- 🔐 **Authentication Middleware**: Simple username-based authentication
- ⚡ **Graceful Shutdown**: Proper cleanup of resources on server shutdown
- 🎨 **Modern UI**: Clean and responsive web interface
- 🚀 **High Performance**: Concurrent client handling with goroutines

## Architecture

The application follows a clean architecture pattern with separation of concerns:

```
chat-app/
├── cmd/
│   └── main.go                 # Application entry point
├── internal/
│   ├── auth/
│   │   └── middleware.go       # Authentication middleware
│   ├── client/
│   │   └── client.go           # WebSocket client management
│   ├── hub/
│   │   └── hub.go              # Central hub for message broadcasting
│   └── message/
│       └── message.go          # Message data structures
├── pkg/
│   └── websocket/
│       └── websocket.go        # WebSocket handler
└── web/
    ├── index.html              # Frontend HTML
    ├── style.css               # Styling
    └── app.js                  # Frontend JavaScript
```

## Components

### Hub

The central message broker that manages:

- Client registration/unregistration
- Room-based message broadcasting
- Client lifecycle management
- Graceful shutdown handling

### Client

Individual WebSocket connection handling:

- Message reading and writing
- Heartbeat/ping-pong mechanism
- Connection cleanup

### WebSocket Handler

HTTP upgrade to WebSocket and connection management

### Auth Middleware

Simple authentication layer to validate user sessions

## Getting Started

### Prerequisites

- Go 1.25.3 or higher
- A modern web browser

### Running the Application

Run the server with:

```bash
go run ./cmd/main.go
```

The server will start on `http://localhost:8080`

## Usage

1. Open your browser and navigate to `http://localhost:8080`
2. Enter a username
3. Choose or create a room (default: "general")
4. Click "Join Chat"
5. Start chatting!

### API Endpoints

- `GET /` - Serves the web interface
- `WebSocket /ws?username=<name>&room=<room>` - WebSocket connection
- `GET /health` - Health check endpoint
- `GET /stats` - Room statistics (returns JSON with active users per room)

## Configuration

### Environment Variables

You can customize the server by modifying `cmd/main.go`:

- `:8080` - Server port (line 58)
- CORS origins - Currently allows all origins (line 50)
- Read/Write/Idle timeouts (lines 60-62)

## Message Types

The application supports several message types:

- `chat` - Regular chat messages
- `join` - User joined notification
- `leave` - User left notification
- `error` - Error messages
- `typing` - Typing indicators

## Security Considerations

⚠️ **Note**: This is a demonstration application. For production use:

1. Implement proper JWT-based authentication
2. Add rate limiting
3. Validate and sanitize all inputs
4. Use secure WebSocket (WSS) in production
5. Implement proper CORS policies
6. Add message encryption for sensitive data
7. Implement proper logging and monitoring

## Technologies Used

- **Go** - Backend programming language
- **Gorilla WebSocket** - WebSocket implementation
- **rs/cors** - CORS middleware
- **Vanilla JavaScript** - Frontend (no frameworks)
- **HTML5 & CSS3** - Modern web interface

## Performance

- Handles multiple concurrent connections efficiently
- Uses buffered channels to prevent blocking
- Implements proper connection pooling
- Graceful degradation on connection errors

## Development

### Project Structure

The project follows Go best practices with:

- Clean separation between handlers, models, and business logic
- Package-level organization
- Context-based cancellation for graceful shutdowns
- Safe concurrent access with mutexes

### Adding Features

To extend the application:

1. Add new message types in `internal/message/message.go`
2. Implement business logic in `internal/hub/hub.go`
3. Add new endpoints in `cmd/main.go`
4. Update frontend in `web/` directory

## Contributing

Contributions are welcome! Please feel free to submit a Pull Request.

## License

This project is open source and available under the MIT License.

## Author

dipendra-mule

---

**Happy Chatting! 💬**

