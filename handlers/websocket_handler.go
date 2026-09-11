package handlers

import (
	"backend/models"
	"context"
	"encoding/json"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func (h *MessageHandler) WebSocket(c *gin.Context) {

	userIDValue, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "user not authenticated",
		})
		return
	}
	userID, ok := userIDValue.(uuid.UUID)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "invalid user id",
		})
		return
	}
	conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		return
	}
	h.Manager.Add(userID, conn)

	ctx := context.Background()

	err = h.RedisService.SetOnline(ctx, userID)

	if err != nil {
		log.Println("Redis presence error:", err)
	}

	defer func() {

		h.Manager.Remove(userID)

		err := h.RedisService.SetOffline(context.Background(), userID)

		if err != nil {
			log.Println("Redis presence delete error:", err)
		}

		err = h.UserService.UpdateLastSeen(userID)

		if err != nil {
			log.Println("Update last seen error:", err)
		}

		conn.Close()
	}()
	for {
		var message models.WebSocketMessage

		err := conn.ReadJSON(&message)

		if err != nil {
			log.Println("WebSocket read error:", err)
			break
		}

		log.Println("Received message:", message)

		switch message.Type {
		case "message":
			createdMessage, err := h.Service.SendMessage(
				userID,
				message,
			)

			if err != nil {
				log.Println("Send message error:", err)

				conn.WriteJSON(gin.H{
					"error": "failed to send message",
				})
				continue
			}

			log.Println("Message saved:", createdMessage)

			event := models.MessageEvent{
				MessageID:  createdMessage.ID,
				SenderID:   userID,
				ReceiverID: message.ReceiverID,
				Content:    message.Content,
			}

			eventData, err := json.Marshal(event)

			if err != nil {
				log.Println("Failed to marshal message event:", err)
				continue
			}

			err = h.RedisService.Publish(
				context.Background(),
				"messages",
				string(eventData),
			)

			if err != nil {
				log.Println("Failed to publish message:", err)
				continue
			}

			log.Println("Message published to Redis:", string(eventData))
		case "read":
			msg, err := h.Service.GetMessageByID(message.MessageID)
			if err != nil {
				log.Println("Failed to get message:", err)
				conn.WriteJSON(gin.H{
					"error": "failed to get message",
				})
				continue
			}
			if msg.ReceiverID != userID {
				log.Println("User is not authorized to read this message")
				conn.WriteJSON(gin.H{
					"error": "user is not authorized to read this message",
				})
				continue
			}
			err = h.Service.UpdateMessageStatus(message.MessageID, userID, models.MessageRead)
			if err != nil {
				log.Println("Failed to mark message as read:", err)
				conn.WriteJSON(gin.H{
					"error": "failed to mark message as read",
				})
				continue

			}
			log.Println("Message marked as read:", message.MessageID)

			event := models.ReadReceiptEvent{
				MessageID: msg.ID,
				SenderID: msg.SenderID,
				ReceiverID: msg.ReceiverID,
			}

			eventData, err := json.Marshal(event)

			if err != nil {
				log.Println("Failed to marshal read receipt:", err)
				continue
			}

			err = h.RedisService.Publish(
				context.Background(),
				"read_receipts",
				string(eventData),
			)

			if err != nil {
				log.Println("Failed to publish read receipt:", err)
				continue
			}

			log.Println("Read receipt published:", string(eventData))

		default:
			log.Println("Unknown message type:", message.Type)
		}

	}

}
