package services

import (
	"backend/models"
	"backend/repositories"

	"github.com/google/uuid"
)

type MessageService struct {
	MessageRepo *repositories.MessageRepository
}

func NewMessageService(
	messageRepo *repositories.MessageRepository,
) *MessageService {
	return &MessageService{
		MessageRepo: messageRepo,
	}
}

func (s *MessageService) SendMessage(senderID uuid.UUID ,req models.WebSocketMessage ) (*models.Message, error) {
	message := models.Message{
		SenderID: 		senderID,
		ReceiverID: 	req.ReceiverID,
		Content: 		req.Content,
	}

	return s.MessageRepo.CreateMessage(message)

}

func (s *MessageService) GetConversation(userID, otherUserID uuid.UUID) ([]models.Message, error) {
	return s.MessageRepo.GetMessages(userID, otherUserID)
}