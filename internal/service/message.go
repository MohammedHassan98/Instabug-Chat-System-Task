package service

import (
	"chat-system/internal/db/models"
	"chat-system/internal/queue"
	"context"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/elastic/go-elasticsearch/v8"
	"github.com/go-redis/redis/v8"
	"gorm.io/gorm"
)

type MessageService struct {
	db    *gorm.DB
	redis *redis.Client
	queue *queue.MessageQueue
	es    *elasticsearch.Client
}

func NewMessageService(db *gorm.DB, redis *redis.Client, queue *queue.MessageQueue, es *elasticsearch.Client) *MessageService {
	return &MessageService{db: db, redis: redis, queue: queue, es: es}
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

func (s *MessageService) SearchMessages(ctx context.Context, chatID uint, query string) ([]models.Message, error) {
	// Construct the search body with partial matching capabilities
	searchBody := map[string]interface{}{
		"query": map[string]interface{}{
			"bool": map[string]interface{}{
				"filter": []interface{}{
					map[string]interface{}{
						"term": map[string]interface{}{
							"ChatID": chatID,
						},
					},
				},
				"must": []interface{}{
					map[string]interface{}{
						"match": map[string]interface{}{
							"Body": map[string]interface{}{
								"query":     query,
								"fuzziness": "AUTO", // Optional: for fuzzy matching
							},
						},
					},
				},
			},
		},
	}

	searchJSON, err := json.Marshal(searchBody)
	if err != nil {
		return nil, fmt.Errorf("error marshaling search body: %w", err)
	}

	res, err := s.es.Search(
		s.es.Search.WithContext(ctx),
		s.es.Search.WithIndex("messages"),
		s.es.Search.WithBody(strings.NewReader(string(searchJSON))))

	if err != nil {
		return nil, fmt.Errorf("error executing search: %w", err)
	}
	defer res.Body.Close()

	var result struct {
		Hits struct {
			Total struct {
				Value    int    `json:"value"`
				Relation string `json:"relation"`
			} `json:"total"`
			Hits []struct {
				Source models.Message `json:"_source"`
			} `json:"hits"`
		} `json:"hits"`
	}

	if err := json.NewDecoder(res.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("error decoding search results: %w", err)
	}

	messages := make([]models.Message, len(result.Hits.Hits))
	for i, hit := range result.Hits.Hits {
		messages[i] = hit.Source
	}

	return messages, nil
}
