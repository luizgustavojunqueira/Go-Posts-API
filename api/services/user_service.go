package services

import (
	"gorm.io/gorm"
	"luizg/PostsAPI/api/models"
)

type UserService struct {
	DB *gorm.DB
}

// Função para salvar um usuário no banco de dados
func (us *UserService) Save(user models.User) (uint, error) {
	result := us.DB.Create(&user)

	// Caso ocorra um erro ao salvar o usuário, retorne o erro
	if result.Error != nil {
		return 0, result.Error
	}

	// Caso contrário, retorne o id do usuário salvo
	return user.ID, nil
}

// Função para buscar um usuário no banco de dados
func (us *UserService) Find(user models.User, id uint) error {
	result := us.DB.Preload("Posts").First(&user, id)

	// Caso ocorra um erro ao buscar o usuário, retorne o erro
	if result.Error != nil {
		return result.Error
	}

	// Caso contrário, retorne nil
	return nil
}

func (us *UserService) FindAll() ([]models.User, error) {
	var users []models.User

	result := us.DB.Preload("Posts").Find(&users)

	if result.Error != nil {
		return nil, result.Error
	}

	return users, nil
}
