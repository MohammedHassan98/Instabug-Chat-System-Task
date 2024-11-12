package models

import (
	"time"
)

type Chat struct {
	ID            uint      `gorm:"primaryKey"`
	ApplicationID uint      `gorm:"not null;index"` // Foreign key referencing Application
	ChatNumber    int       `gorm:"not null"`
	MessagesCount int       `gorm:"default:0"`
	Messages      []Message `gorm:"foreignKey:ChatID"`
	CreatedAt     time.Time
	UpdatedAt     time.Time

	Application Application `gorm:"constraint:OnDelete:CASCADE"`

	// Add unique composite index for ApplicationID and ChatNumber
	// to prevent duplicate chat numbers within an application
	Index []string `gorm:"uniqueIndex:idx_application_chat_number"`
}
