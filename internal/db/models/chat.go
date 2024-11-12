package models

import (
	"time"
)

type Chat struct {
	ID            uint      `gorm:"primaryKey"`
	ApplicationID uint      `gorm:"not null;index"`
	ChatNumber    int       `gorm:"not null;uniqueIndex:idx_application_chat_number"`
	MessagesCount int       `gorm:"default:0"`
	Messages      []Message `gorm:"foreignKey:ChatID"`
	CreatedAt     time.Time
	UpdatedAt     time.Time

	Application Application `gorm:"constraint:OnDelete:CASCADE"`
}
