package services

import (
	"backend/models"
	"backend/repositories"
	"errors"

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

func (s *MessageService) SendMessage(senderID uuid.UUID, req models.WebSocketMessage) (*models.Message, error) {
	message := models.Message{
		SenderID:   senderID,
		ReceiverID: req.ReceiverID,
		Content:    req.Content,
	}

	return s.MessageRepo.SendMessage(message)

}

func (s *MessageService) GetConversation(userID, otherUserID uuid.UUID) ([]models.Message, error) {
	return s.MessageRepo.GetMessages(userID, otherUserID)
}

func (s *MessageService) UpdateMessageStatus(messageID uuid.UUID, userID uuid.UUID, status string) error {
	if status != models.MessageDelivered &&
		status != models.MessageRead {
		return errors.New("invalid message status")
	}
	return s.MessageRepo.UpdateMessageStatus(messageID, userID, status)
}

func (s *MessageService) GetMessageByID(messageID uuid.UUID) (*models.Message, error) {
   return s.MessageRepo.GetMessageByID(messageID)
}
