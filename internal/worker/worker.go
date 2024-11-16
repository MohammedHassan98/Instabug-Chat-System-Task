package worker

import (
	"bytes"
	"chat-system/internal/db"
	"chat-system/internal/db/models"
	"context"
	"encoding/json"
	"fmt"
	"log"
	"sync"
	"time"

	"chat-system/internal/queue"

	"github.com/elastic/go-elasticsearch/v8"
)

type Worker struct {
	queue *queue.MessageQueue
	es    *elasticsearch.Client
	mu    sync.Mutex // Mutex to protect critical sections
}

func NewWorker(queue *queue.MessageQueue, es *elasticsearch.Client) *Worker {
	return &Worker{queue: queue, es: es}
}

func (w *Worker) Start(ctx context.Context) {
	go w.processQueue(ctx)
	go w.updateCounters(ctx)
}

func (w *Worker) processQueue(ctx context.Context) {
	for {
		select {
		case <-ctx.Done():
			return
		default:
			result, err := db.Redis.BRPop(ctx, 0, "message_queue").Result()
			if err != nil {
				log.Printf("Error popping from queue: %v", err)
				continue
			}

			var qm queue.QueuedMessage
			if err := json.Unmarshal([]byte(result[1]), &qm); err != nil {
				log.Printf("Error unmarshaling message: %v", err)
				continue
			}

			switch qm.Type {
			case "chat_creation":
				w.processChatCreation(ctx, qm.Payload)
			case "message_creation":
				w.processMessageCreation(ctx, qm.Payload)
			}
		}
	}
}

func (w *Worker) processChatCreation(ctx context.Context, payload json.RawMessage) {
	var data struct {
		AppID      uint `json:"app_id"`
		ChatNumber int  `json:"chat_number"`
	}

	if err := json.Unmarshal(payload, &data); err != nil {
		log.Printf("Error unmarshaling chat creation payload: %v", err)
		return
	}

	w.mu.Lock() // Lock before creating chat
	defer w.mu.Unlock()

	chat := &models.Chat{
		ApplicationID: data.AppID,
		ChatNumber:    data.ChatNumber,
	}

	if err := db.GormDB.Create(chat).Error; err != nil {
		log.Printf("Error creating chat: %v", err)
		return
	}
}

func (w *Worker) processMessageCreation(ctx context.Context, payload json.RawMessage) {
	var data struct {
		ChatID        uint   `json:"chat_id"`
		MessageNumber int    `json:"message_number"`
		Body          string `json:"body"`
	}

	if err := json.Unmarshal(payload, &data); err != nil {
		log.Printf("Error unmarshaling message creation payload: %v", err)
		return
	}

	w.mu.Lock() // Lock before creating chat
	defer w.mu.Unlock()

	message := &models.Message{
		ChatID:        data.ChatID,
		MessageNumber: data.MessageNumber,
		Body:          data.Body,
	}

	if err := db.GormDB.Create(message).Error; err != nil {
		log.Printf("Error creating message: %v", err)
		return
	}

	// Index the message in Elasticsearch
	messageJSON, err := json.Marshal(message)
	if err != nil {
		log.Printf("Error marshaling message for Elasticsearch: %v", err)
		return
	}

	_, err = w.es.Index(
		"messages", // Index name
		bytes.NewReader(messageJSON),
		w.es.Index.WithDocumentID(fmt.Sprintf("%d-%d", message.ChatID, message.MessageNumber)),
		w.es.Index.WithRefresh("true"), // Refresh the index
	)
	if err != nil {
		log.Printf("Error indexing message in Elasticsearch: %v", err)
		return
	}
	if err != nil {
		log.Printf("Error indexing message in Elasticsearch: %v", err)
		return
	}
}

func (w *Worker) updateCounters(ctx context.Context) {
	ticker := time.NewTicker(15 * time.Minute)
	for {
		select {
		case <-ticker.C:
			w.syncCounters(ctx)
		case <-ctx.Done():
			return
		}
	}
}

func (w *Worker) syncCounters(ctx context.Context) {
	// Update application chat counts
	db.GormDB.Exec(`
		UPDATE applications a
		SET chats_count = (
			SELECT COUNT(*) FROM chats
			WHERE application_id = a.id
		)
	`)

	// Update chat message counts
	db.GormDB.Exec(`
		UPDATE chats c
		SET messages_count = (
			SELECT COUNT(*) FROM messages
			WHERE chat_id = c.id
		)
	`)
}
