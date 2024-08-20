package services

import (
	"gorm.io/gorm"

	"luizg/PostsAPI/api/models"
)

type PostService struct {
	DB *gorm.DB
}

func (ps *PostService) Save(post models.Post) (uint, error) {
	result := ps.DB.Create(&post)

	if result.Error != nil {
		return 0, result.Error
	}

	return post.ID, nil
}

func (ps *PostService) Find(post models.Post, id uint) error {
	result := ps.DB.Preload("User").First(&post, id)

	if result.Error != nil {
		return result.Error
	}

	return nil
}

func (ps *PostService) FindAll() ([]models.Post, error) {
	var posts []models.Post

	result := ps.DB.Find(&posts)

	if result.Error != nil {
		return nil, result.Error
	}

	return posts, nil
}
