package handlers

import (
	"backend/models"
	"context"
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

	err = h.Redis.Set(
		ctx,
		"user:"+userID.String(),
		"online",
		0,
	).Err()

	if err != nil {
		log.Println("Redis presence error:", err)
	}

	defer func() {

		h.Manager.Remove(userID)

		err := h.Redis.Del(
			context.Background(),
			"user:"+userID.String(),
		).Err()

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

		receiverConn, exists := h.Manager.Get(message.ReceiverID)

		if !exists {
			log.Println("Receiver is offline:", message.ReceiverID)
			continue
		}

		err = receiverConn.WriteJSON(createdMessage)

		if err != nil {
			log.Println("WebSocket write error:", err)
			h.Manager.Remove(message.ReceiverID)
		}
	}

}
