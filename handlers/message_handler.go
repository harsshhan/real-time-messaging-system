package handlers

import (
	"backend/models"
	"backend/services"
	"backend/ws"
	"context"
	"encoding/json"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type MessageHandler struct {
	Service      *services.MessageService
	Manager      *ws.ConnectionManager
	RedisService *services.RedisService
	UserService  *services.UserService
}

func NewMessageHandler(service *services.MessageService, manager *ws.ConnectionManager, RedisService *services.RedisService, userService *services.UserService) *MessageHandler {
	return &MessageHandler{
		Service:      service,
		Manager:      manager,
		RedisService: RedisService,
		UserService:  userService,
	}

}

func (h *MessageHandler) SendMessage(c *gin.Context) {

	var req models.WebSocketMessage

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "invalid request",
		})
		return
	}

	userID, exists := c.Get("userID")

	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "user not authenticated",
		})
		return
	}

	senderID, ok := userID.(uuid.UUID)

	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "invalid user id",
		})
		return
	}

	message, err := h.Service.SendMessage(senderID, req)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "failed to send message",
		})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": message,
	})
}

func (h *MessageHandler) GetConversation(c *gin.Context) {

	userIDstr, exists := c.Get("userID")

	if !exists {
		log.Println("[GetConversation] Error: user not authenticated in context")
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "user not authenticated",
		})
		return
	}

	userID, ok := userIDstr.(uuid.UUID)

	if !ok {
		log.Println("[GetConversation] Error: invalid user id in context")
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "invalid user id",
		})
		return
	}

	paramOtherID := c.Param("otherUserID")
	if paramOtherID == "" {
		paramOtherID = c.Param("userId")
	}

	otherUserID, err := uuid.Parse(paramOtherID)

	if err != nil {
		log.Printf("[GetConversation] Error parsing otherUserID ('%s'): %v\n", paramOtherID, err)
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "invalid user id",
		})
		return
	}

	messages, err := h.Service.GetConversation(
		userID,
		otherUserID,
	)

	if err != nil {
		log.Printf("[GetConversation] Service error: %v\n", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "failed to get messages",
		})
		return
	}

	c.JSON(http.StatusOK, messages)
}

func (h *MessageHandler) ListenForMessages(ctx context.Context) {

	pubsub := h.RedisService.Subscribe(ctx, "messages", "read_receipts")
	defer pubsub.Close()

	for {
		message, err := pubsub.ReceiveMessage(ctx)

		if err != nil {
			log.Println("Redis subscriber error:", err)
			return
		}

		switch message.Channel {
		case "messages":

			var event models.MessageEvent

			err = json.Unmarshal([]byte(message.Payload), &event)
			if err != nil {
				log.Println("Failed to unmarshal Redis message:", err)
				continue
			}

			log.Println("Redis event:", event)

			receiverConn, exists := h.Manager.Get(event.ReceiverID)
			if !exists {
				log.Println("Receiver is offline:", event.ReceiverID)
				continue
			}
			err = receiverConn.WriteJSON(event)
			if err != nil {
				log.Println("WebSocket write error:", err)
				h.Manager.Remove(event.ReceiverID)
				continue
			}
			err = h.Service.UpdateMessageStatus(event.MessageID, event.ReceiverID, models.MessageDelivered)

			if err != nil {
				log.Println("Failed to mark message as delivered:", err)
				continue
			}
			log.Println("Message delivered:", event.MessageID)

		case "read_receipts":
			var readReceiptEvent models.ReadReceiptEvent

			err := json.Unmarshal(
				[]byte(message.Payload),
				&readReceiptEvent,
			)

			if err != nil {
				log.Println("Failed to unmarshal read receipt:", err)
				continue
			}

			log.Println("Redis read receipt:", readReceiptEvent)

			senderConn, exists := h.Manager.Get(readReceiptEvent.SenderID)

			if !exists {
				log.Println("Sender is offline:", readReceiptEvent.SenderID)
				continue
			}

			err = senderConn.WriteJSON(readReceiptEvent)

			if err != nil {
				log.Println("WebSocket write error:", err)
				h.Manager.Remove(readReceiptEvent.SenderID)
				continue
			}

			log.Println("Read receipt delivered:", readReceiptEvent.MessageID)
		}
	}
}
