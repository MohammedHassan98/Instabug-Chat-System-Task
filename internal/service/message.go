package service

import (
	"chat-system/internal/db/models"
	"chat-system/internal/queue"
	"context"
	"fmt"

	"github.com/go-redis/redis/v8"
	"gorm.io/gorm"
)

type MessageService struct {
	db    *gorm.DB
	redis *redis.Client
	queue *queue.MessageQueue
}

func NewMessageService(db *gorm.DB, redis *redis.Client, queue *queue.MessageQueue) *MessageService {
	return &MessageService{db: db, redis: redis, queue: queue}
}

func (s *MessageService) CreateMessage(ctx context.Context, chatID uint, body string) (*models.Message, error) {
	// Get next message number using Redis
	msgNumKey := fmt.Sprintf("chat:%d:next_msg_num", chatID)
	msgNum, err := s.redis.Incr(ctx, msgNumKey).Result()
	if err != nil {
		return nil, err
	}

	message := &models.Message{
		ChatID:        chatID,
		MessageNumber: int(msgNum),
		Body:          body,
	}

	// Queue the message creation
	payload := struct {
		ChatID        uint   `json:"chat_id"`
		MessageNumber int    `json:"message_number"`
		Body          string `json:"body"`
	}{
		ChatID:        chatID,
		MessageNumber: int(msgNum),
		Body:          body,
	}

	if err := s.queue.Enqueue(ctx, "message_creation", payload); err != nil {
		return nil, err
	}

	return message, nil
}

func (s *MessageService) GetMessagesByChatNumberAndToken(ctx context.Context, token string, chatNumber uint) ([]models.Message, error) {
	var messages []models.Message

	// Join with the Application model to filter by token
	if err := s.db.Table("messages").
		Select("messages.*").
		Joins("JOIN chats ON chats.id = messages.chat_id").
		Joins("JOIN applications ON applications.id = chats.application_id").
		Where("applications.token = ? AND chats.chat_number = ?", token, chatNumber).
		Find(&messages).Error; err != nil {
		return nil, err
	}

	return messages, nil
}
