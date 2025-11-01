package hub

import (
	"context"
	"log"
	"sync"
	"time"

	"github.com/dipendra-mule/chat-app/internal/client"
	"github.com/dipendra-mule/chat-app/internal/message"
)

type Hub struct {
	mu sync.RWMutex

	// Registered cleints by room
	rooms map[string]map[*client.Client]bool

	// Inbound messages from clients
	broadcast chan *message.Message

	// register req
	register chan *client.Client

	// unregister req
	unregister chan *client.Client

	// context for graceful shutdown
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
	log.Println("chat hub started")

	for {
		select {
		case client := <-h.register:
			h.mu.Lock()
			if h.rooms[client.Room] == nil {
				h.rooms[client.Room] = make(map[*client.Client]bool)
			}
			h.rooms[client.Room][client] = true
			h.mu.Unlock()

			// notifyt room about user

			joinMsg := &message.Message{
				Type:      message.MessageTypeJoin,
				Username:  client.Username,
				Content:   "joined the room",
				Timestamp: time.Now(),
				Room:      client.Room,
			}
			h.broadcastMessage(joinMsg)
			log.Printf("Client %s joined the room %s", client.Username, client.Room)
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
					// clean up empty room
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
			// close all client connections
			for _, room := range h.rooms {
				for client := range room {
					close(client.Send)
				}
			}
			h.mu.Unlock()
			log.Println("chat hub stopped")
			return
		}

	}
}
