package models

import (
	"time"
)

type Chat struct {
	ID            uint      `gorm:"primaryKey"`
	ApplicationID uint      `gorm:"not null;index"`
	ChatNumber    int       `gorm:"not null;index:idx_application_chat_number"`
	MessagesCount int       `gorm:"default:0"`
	Messages      []Message `gorm:"foreignKey:ChatID"`
	CreatedAt     time.Time
	UpdatedAt     time.Time

	CompositeIndex string      `gorm:"index:idx_application_chat_number_application_id,unique;not null"`
	Application    Application `gorm:"constraint:OnDelete:CASCADE"`
}
