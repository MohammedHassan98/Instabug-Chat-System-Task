package models

import (
	"time"
)

type Message struct {
	ID            uint   `gorm:"primaryKey"`
	ChatID        uint   `gorm:"not null;index"` // Foreign key referencing Chat
	MessageNumber int    `gorm:"not null;index"`
	Body          string `gorm:"type:text;not null"`
	CreatedAt     time.Time
	UpdatedAt     time.Time

	Chat Chat `gorm:"constraint:OnDelete:CASCADE"`
}
