package hub

import (
	"context"
	"log"
	"sync"
	"time"

	"github.com/dipendra-mule/chat-app/internal/client"  // client package
	"github.com/dipendra-mule/chat-app/internal/message" // message package
)

type Hub struct {
	mu sync.RWMutex

	// Registered clients by room
	rooms map[string]map[*client.Client]bool

	// Inbound messages from clients
	broadcast chan *message.Message

	// Register requests
	register chan *client.Client

	// Unregister requests
	unregister chan *client.Client

	// Context for graceful shutdown
	ctx    context.Context
	cancel context.CancelFunc
}

func NewHub() *Hub {
	ctx, cancel := context.WithCancel(context.Background())
	return &Hub{
		rooms:      make(map[string]map[*client.Client]bool),
		broadcast:  make(chan *message.Message),
		register:   make(chan *client.Client),
		unregister: make(chan *client.Client),
		ctx:        ctx,
		cancel:     cancel,
	}
}

func (h *Hub) Run() {
	log.Println("Chat hub started")

	for {
		select {
		case client := <-h.register:
			h.mu.Lock()
			if h.rooms[client.Room] == nil {
				h.rooms[client.Room] = make(map[*client.Client]bool)
			}
			h.rooms[client.Room][client] = true
			h.mu.Unlock()

			// Notify room about new user
			joinMsg := &message.Message{
				Type:      message.MessageTypeJoin,
				Username:  client.Username,
				Content:   "joined the room",
				Timestamp: time.Now(),
				Room:      client.Room,
			}
			h.broadcastMessage(joinMsg)

			log.Printf("Client %s joined room %s", client.Username, client.Room)

		case client := <-h.unregister:
			h.mu.Lock()
			if room, ok := h.rooms[client.Room]; ok {
				if _, ok := room[client]; ok {
					delete(room, client)
					close(client.Send)

					// Notify room about user leaving
					leaveMsg := &message.Message{
						Type:      message.MessageTypeLeave,
						Username:  client.Username,
						Content:   "left the room",
						Timestamp: time.Now(),
						Room:      client.Room,
					}
					go func() {
						h.broadcast <- leaveMsg
					}()

					// Clean up empty rooms
					if len(room) == 0 {
						delete(h.rooms, client.Room)
					}
				}
			}
			h.mu.Unlock()

		case message := <-h.broadcast:
			h.broadcastMessage(message)

		case <-h.ctx.Done():
			h.mu.Lock()
			// Close all client connections
			for _, room := range h.rooms {
				for client := range room {
					close(client.Send)
				}
			}
			h.mu.Unlock()
			log.Println("Chat hub stopped")
			return
		}
	}
}

func (h *Hub) broadcastMessage(msg *message.Message) {
	h.mu.RLock()
	defer h.mu.RUnlock()

	if room, ok := h.rooms[msg.Room]; ok {
		for client := range room {
			select {
			case client.Send <- msg:
			default:
				close(client.Send)
				delete(room, client)
			}
		}
	}
}

func (h *Hub) GetRoomStats() map[string]int {
	h.mu.RLock()
	defer h.mu.RUnlock()

	stats := make(map[string]int)
	for room, clients := range h.rooms {
		stats[room] = len(clients)
	}
	return stats
}

func (h *Hub) Stop() {
	h.cancel()
}
