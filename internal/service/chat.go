package service

import (
	"chat-system/internal/db/models"
	"chat-system/internal/queue"
	"context"
	"fmt"

	"github.com/go-redis/redis/v8"
	"gorm.io/gorm"
)

type ChatService struct {
	db    *gorm.DB
	redis *redis.Client
	queue *queue.MessageQueue
}

func NewChatService(db *gorm.DB, redis *redis.Client, queue *queue.MessageQueue) *ChatService {
	return &ChatService{db: db, redis: redis, queue: queue}
}

func (s *ChatService) CreateChat(ctx context.Context, appToken string) (*models.Chat, error) {
	// Start transaction
	tx := s.db.Begin()
	if tx.Error != nil {
		return nil, tx.Error
	}

	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	// Get application
	var app models.Application
	if err := tx.Where("token = ?", appToken).First(&app).Error; err != nil {
		tx.Rollback()
		return nil, err
	}

	// Get next chat number using Redis
	chatNumKey := fmt.Sprintf("app:%d:next_chat_num", app.ID)
	chatNum, err := s.redis.Incr(ctx, chatNumKey).Result()
	if err != nil {
		tx.Rollback()
		return nil, err
	}

	chat := &models.Chat{
		ApplicationID: app.ID,
		ChatNumber:    int(chatNum),
	}

	// Queue the chat creation
	payload := struct {
		AppID      uint `json:"app_id"`
		ChatNumber int  `json:"chat_number"`
	}{
		AppID:      app.ID,
		ChatNumber: int(chatNum),
	}

	if err := s.queue.Enqueue(ctx, "chat_creation", payload); err != nil {
		tx.Rollback()
		return nil, err
	}

	if err := tx.Commit().Error; err != nil {
		return nil, err
	}

	return chat, nil
}
