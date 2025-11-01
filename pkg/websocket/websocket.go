package websocket

import (
	"context"
	"log"
	"net/http"
	"time"

	"github.com/dipendra-mule/chat-app/internal/client" // client package
	// message package
	"github.com/dipendra-mule/chat-app/internal/hub" // hub package

	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true // In production, validate origins
	},
}

type Handler struct {
	hub *hub.Hub
}

func NewHandler(hub *hub.Hub) *Handler {
	return &Handler{hub: hub}
}

func (h *Handler) ServeWebSocket(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Printf("WebSocket upgrade error: %v", err)
		return
	}
	defer conn.Close()

	// Extract user info (in real app, use JWT or session)
	username := r.URL.Query().Get("username")
	room := r.URL.Query().Get("room")

	if username == "" || room == "" {
		conn.WriteJSON(map[string]string{"error": "username and room are required"})
		return
	}

	client := client.NewClient(generateID(), username, room, conn)

	// Register client
	h.hub.Register <- client

	// Create context for this connection
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Start goroutines for reading and writing
	go client.WritePump(ctx)
	client.ReadPump(ctx, h.hub.Broadcast)

	// Unregister client
	h.hub.Unregister <- client
}

func generateID() string {
	return string(time.Now().UnixNano())
}
