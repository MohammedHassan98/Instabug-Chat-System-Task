package models

import (
	"time"
)

type Message struct {
	ID            uint   `gorm:"primaryKey"`
	ChatID        uint   `gorm:"not null;index"`
	MessageNumber int    `gorm:"not null;index:idx_chat_message_number"`
	Body          string `gorm:"type:text;not null"`
	CreatedAt     time.Time
	UpdatedAt     time.Time

	Chat Chat `gorm:"constraint:OnDelete:CASCADE"`
}
