package models

import (
	"gorm.io/gorm"
	"time"
)

type Post struct {
	ID           uint      `gorm:"primaryKey" json:"post_id"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
	Title        string    `json:"title"`
	Content      string    `json:"content"`
	UserID       uint      `json:"user_id"`
	UserFullName string    `json:"user_full_name"`
}

type CreatePost struct {
	Title   string `json:"title" binding:"required"`
	Content string `json:"content" binding:"required"`
}

type PostService struct {
	DB *gorm.DB
}

// Function to save a post in the database
func (ps *PostService) Save(post Post) (uint, error) {
	result := ps.DB.Create(&post)

	if result.Error != nil {
		return 0, result.Error
	}

	return post.ID, nil
}

// Function to find a post in the database
func (ps *PostService) Find(post Post, id uint) error {
	result := ps.DB.Preload("User").First(&post, id)

	if result.Error != nil {
		return result.Error
	}

	return nil
}

// Function to find all posts in the database
func (ps *PostService) FindAll() ([]Post, error) {
	var posts []Post

	result := ps.DB.Find(&posts)

	if result.Error != nil {
		return nil, result.Error
	}

	return posts, nil
}

// Function to find a post by ID in the database
func (ps *PostService) FindByID(id uint) (Post, error) {
	var post Post

	result := ps.DB.First(&post, id)

	if result.Error != nil {
		return post, result.Error
	}

	return post, nil
}

// Function to find all posts by user ID in the database
func (ps *PostService) FindByUserID(userID uint) ([]Post, error) {
	var posts []Post

	result := ps.DB.Where("user_id = ?", userID).Find(&posts)

	if result.Error != nil {
		return nil, result.Error
	}

	return posts, nil
}

// Function to delete a post in the database
func (ps *PostService) Delete(post Post) (uint, error) {
	result := ps.DB.Delete(&post)

	if result.Error != nil {
		return 0, result.Error
	}

	return post.ID, nil
}
