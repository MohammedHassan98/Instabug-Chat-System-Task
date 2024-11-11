package models

import (
	"time"
)

type Chat struct {
	ID            uint      `gorm:"primaryKey"`
	ApplicationID uint      `gorm:"not null;index"` // Foreign key referencing Application
	ChatNumber    int       `gorm:"not null;index"`
	MessagesCount int       `gorm:"default:0"`
	Messages      []Message `gorm:"foreignKey:ChatID"`
	CreatedAt     time.Time
	UpdatedAt     time.Time

	Application Application `gorm:"constraint:OnDelete:CASCADE"`
}
