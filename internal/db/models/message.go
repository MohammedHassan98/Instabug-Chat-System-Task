package models

import (
	"time"
)

type Message struct {
	ID            uint   `gorm:"primaryKey"`
	ChatID        uint   `gorm:"not null;index"`
	MessageNumber int    `gorm:"not null"`
	Body          string `gorm:"type:text;not null"`
	CreatedAt     time.Time
	UpdatedAt     time.Time

	Chat Chat `gorm:"constraint:OnDelete:CASCADE"`

	// Add unique composite index for ChatID and MessageNumber
	Index []string `gorm:"uniqueIndex:idx_chat_message_number"`
}
