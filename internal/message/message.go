package message

import (
	"encoding/json"
	"time"
)

type MessageType string

const (
	MessageTypeChat   MessageType = "chat"
	MessageTypeJoin   MessageType = "join"
	MessageTypeLeave  MessageType = "leave"
	MessageTypeError  MessageType = "error"
	MessageTypeTyping MessageType = "typing"
)

type Message struct {
	ID        string      `json:"id, omitempty"`
	Type      MessageType `json:"type, omitempty"`
	Username  string      `json:"username, omitempty"`
	Content   string      `json:"content, omitempty"`
	Timestamp time.Time   `json:"timestamp, omitempty"`
	Room      string      `json:"room, omitempty"`
}

func (m *Message) ToJSON() []byte {
	bytes, _ := json.Marshal(m)
	return bytes
}

func FromJSON(data []byte) (*Message, error) {
	var msg Message
	err := json.Unmarshal(data, &msg)
	return &msg, err
}
