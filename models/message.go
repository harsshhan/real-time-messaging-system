package models

import (
	"time"

	"github.com/google/uuid"
)

type Message struct {
	ID uuid.UUID `json:"id"`
	SenderID uuid.UUID `json:"sender_id"`
	ReceiverID uuid.UUID `json:"receiver_id"`
	Content string `json:"content"`
	CreatedAt time.Time `json:"created_at"`
}

type WebSocketMessage struct {
	ReceiverID uuid.UUID `json:"receiver_id"`
	Content string `json:"content"`
}

type MessageResponse struct {
	ID         uuid.UUID `json:"id"`
	SenderID   uuid.UUID `json:"sender_id"`
	ReceiverID uuid.UUID `json:"receiver_id"`
	Content    string    `json:"content"`
	CreatedAt  time.Time `json:"created_at"`
}

