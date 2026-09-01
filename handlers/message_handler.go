package handlers

import (
	"backend/models"
	"backend/services"
	"backend/ws"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type MessageHandler struct{
	Service *services.MessageService
	Manager *ws.ConnectionManager
}

func NewMessageHandler(service *services.MessageService	,manager *ws.ConnectionManager) *MessageHandler {
	return &MessageHandler{
		Service: service,
		Manager: manager,
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