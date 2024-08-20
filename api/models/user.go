package models

import (
	"time"
)

type User struct {
	ID        uint      `gorm:"primaryKey"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	FirstName string    `json:"firstname"`
	LastName  string    `json:"lastname"`
	Email     string    `json:"email"`
	Password  string    `json:"password,omitempty"`
	Posts     []Post    `json:"posts" gorm:"constraint:OnDelete:CASCADE"`
}
