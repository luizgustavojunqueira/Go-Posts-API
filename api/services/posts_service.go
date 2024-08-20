package services

import (
	"gorm.io/gorm"

	"luizg/PostsAPI/api/models"
)

type PostService struct {
	DB *gorm.DB
}

// Function to save a post in the database
func (ps *PostService) Save(post models.Post) (uint, error) {
	result := ps.DB.Create(&post)

	if result.Error != nil {
		return 0, result.Error
	}

	return post.ID, nil
}

// Function to find a post in the database
func (ps *PostService) Find(post models.Post, id uint) error {
	result := ps.DB.Preload("User").First(&post, id)

	if result.Error != nil {
		return result.Error
	}

	return nil
}

// Function to find all posts in the database
func (ps *PostService) FindAll() ([]models.Post, error) {
	var posts []models.Post

	result := ps.DB.Find(&posts)

	if result.Error != nil {
		return nil, result.Error
	}

	return posts, nil
}

// Function to find a post by ID in the database
func (ps *PostService) FindByID(id uint) (models.Post, error) {
	var post models.Post

	result := ps.DB.First(&post, id)

	if result.Error != nil {
		return post, result.Error
	}

	return post, nil
}
