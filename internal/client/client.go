package client

import (
	"context"
	"log"
	"time"

	"github.com/dipendra-mule/chat-app/internal/message"
	"github.com/gorilla/websocket"
)

type Client struct {
	ID       string
	Username string
	Room     string
	conn     *websocket.Conn
	send     chan *message.Message
	// mu       sync.Mutex
}

func NewClient(id, username, room string, conn *websocket.Conn) *Client {
	return &Client{
		ID:       id,
		Username: username,
		Room:     room,
		conn:     conn,
		send:     make(chan *message.Message, 256),
	}
}

func (c *Client) WritePump(ctx context.Context) {
	ticker := time.NewTicker(54 * time.Second)
	defer func() {
		ticker.Stop()
		c.conn.Close()
	}()

	for {
		select {
		case message, ok := <-c.send:
			c.conn.SetWriteDeadline(time.Now().Add(10 * time.Second))
			if !ok {
				c.conn.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}

			w, err := c.conn.NextWriter(websocket.TextMessage)
			if err != nil {
				return
			}
			w.Write(message.ToJSON())

			// Add queued messages to the current websocket message
			n := len(c.send)
			for i := 0; i < n; i++ {
				w.Write([]byte("\n"))
				w.Write((<-c.send).ToJSON())
			}

			if err := w.Close(); err != nil {
				return
			}

		case <-ticker.C:
			c.conn.SetWriteDeadline(time.Now().Add(10 * time.Second))
			if err := c.conn.WriteMessage(websocket.PingMessage, nil); err != nil {
				return
			}

		case <-ctx.Done():
			return
		}
	}
}

func (c *Client) ReadPump(ctx context.Context, broadcast chan<- *message.Message) {
	defer func() {
		c.conn.Close()
	}()

	c.conn.SetReadLimit(5120)
	c.conn.SetReadDeadline(time.Now().Add(60 * time.Second))
	c.conn.SetPongHandler(func(string) error {
		c.conn.SetReadDeadline(time.Now().Add(60 * time.Second))
		return nil
	})

	for {
		select {
		case <-ctx.Done():
			return
		default:
			_, raw, err := c.conn.ReadMessage()
			if err != nil {
				if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
					log.Printf("error: %v", err)
				}
				return
			}

			msg, err := message.FromJSON(raw)
			if err != nil {
				log.Printf("error parsing message: %v", err)
				continue
			}

			msg.Username = c.Username
			msg.Room = c.Room
			msg.Timestamp = time.Now()

			broadcast <- msg
		}
	}
}

func (c *Client) Send(msg *message.Message) {
	select {
	case c.send <- msg:
	default:
		close(c.send)
	}
}
