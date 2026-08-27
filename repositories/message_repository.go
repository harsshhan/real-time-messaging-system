package repositories

import (
	"backend/models"
	"context"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
)

type MessageRepository struct {
	DB *pgxpool.Pool
}

func NewMessageRepository(db *pgxpool.Pool) *MessageRepository {
	return &MessageRepository{
		DB: db,
	}
}

func (r *MessageRepository) CreateMessage(message models.Message) (*models.Message,error) {
	query := `insert into messages (sender_id,receiver_id,content) values ($1,$2,$3) RETURNING id, created_at`
	
	var createdMessage models.Message
	
	err := r.DB.QueryRow(context.Background(),query,message.SenderID,message.ReceiverID,message.Content).Scan(&createdMessage.ID,&createdMessage.CreatedAt)

	if(err!=nil){
		return nil,err
	}

	createdMessage.SenderID = message.SenderID
	createdMessage.ReceiverID = message.ReceiverID
	createdMessage.Content = message.Content

	return &createdMessage,nil

}

func (r *MessageRepository) GetMessages(user1 uuid.UUID, user2 uuid.UUID) ([]models.Message,error){
	query := `Select id,sender_id,receiver_id,content,created_at from messages where 
	(sender_id = $1 AND receiver_id = $2) OR (sender_id = $2 AND receiver_id = $1)
	ORDER BY created_at ASC` 

	rows,err := r.DB.Query(context.Background(),query,user1,user2)

	if(err!=nil){
		return nil,err
	}

	defer rows.Close()

	var messages []models.Message

	for rows.Next(){
		var message models.Message

		err := rows.Scan(&message.ID, &message.SenderID, &message.ReceiverID, &message.Content, &message.CreatedAt)
		if err != nil {
			return nil, err
		}
		messages = append(messages, message)
	}

	return messages, nil
}