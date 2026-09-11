package models

import (
	"time"

	"github.com/google/uuid"
)


const (
	MessageSent      = "sent"
	MessageDelivered = "delivered"
	MessageRead      = "read"
)

type Message struct {
	ID uuid.UUID `json:"id"`
	SenderID uuid.UUID `json:"sender_id"`
	ReceiverID uuid.UUID `json:"receiver_id"`
	Content string `json:"content"`
	Status string `json:"status"`
	CreatedAt time.Time `json:"created_at"`
}

type WebSocketMessage struct {
	Type       string    `json:"type"`
	MessageID  uuid.UUID `json:"message_id,omitempty"`
	ReceiverID uuid.UUID `json:"receiver_id,omitempty"`
	Content    string    `json:"content,omitempty"`
}

type MessageResponse struct {
	ID         uuid.UUID `json:"id"`
	SenderID   uuid.UUID `json:"sender_id"`
	ReceiverID uuid.UUID `json:"receiver_id"`
	Content    string    `json:"content"`
	Status 	   string 	 `json:"status"`
	CreatedAt  time.Time `json:"created_at"`
}

type MessageEvent struct {

	MessageID  uuid.UUID `json:"message_id"`
	SenderID   uuid.UUID `json:"sender_id"`
	ReceiverID uuid.UUID `json:"receiver_id"`
	Content    string    `json:"content"`

}
