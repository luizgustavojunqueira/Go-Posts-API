package models

import (
	"gorm.io/gorm"
	"time"
)

type Post struct {
	ID        uint           `gorm:"primaryKey"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"deleted_at,omitempty"`
	Title     string         `json:"title"`
	Content   string         `json:"content"`
	UserID    int            `json:"authorId"`
}
