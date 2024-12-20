package models

import (
	"time"

	"gorm.io/gorm"
)

type Application struct {
	ID         uint   `gorm:"primarykey"`
	Name       string `gorm:"size:255;not null"`
	Token      string `gorm:"size:255;not null;unique;index"`
	ChatsCount int    `gorm:"default:0"`
	Chats      []Chat `gorm:"foreignKey:ApplicationID"`

	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt `gorm:"index"`
}
