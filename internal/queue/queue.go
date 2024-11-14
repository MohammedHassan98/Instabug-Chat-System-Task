package queue

import (
	"context"
	"encoding/json"

	"github.com/go-redis/redis/v8"
)

type MessageQueue struct {
	redis *redis.Client
}

type QueuedMessage struct {
	Type    string          `json:"type"`
	Payload json.RawMessage `json:"payload"`
}

func NewMessageQueue(redis *redis.Client) *MessageQueue {
	return &MessageQueue{redis: redis}
}

func (mq *MessageQueue) Enqueue(ctx context.Context, msgType string, payload interface{}) error {
	data, err := json.Marshal(payload)
	if err != nil {
		return err
	}

	qm := QueuedMessage{
		Type:    msgType,
		Payload: data,
	}

	qmJSON, err := json.Marshal(qm)
	if err != nil {
		return err
	}

	return mq.redis.LPush(ctx, "message_queue", qmJSON).Err()
}
