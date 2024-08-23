package services

import (
	"gorm.io/gorm"
	"luizg/PostsAPI/api/models"
)

type UserService struct {
	DB *gorm.DB
}

// Function to save a user in the database
func (us *UserService) Save(user models.User) (uint, error) {
	result := us.DB.Create(&user)

	// If an error occurs while saving the user, return the error
	if result.Error != nil {
		return 0, result.Error
	}

	// Otherwise, return the new user ID
	return user.ID, nil
}

// Function to delete a user by email in the database
func (us *UserService) Delete(email string) error {
	result := us.DB.Where("email = ?", email).Delete(&models.User{})

	if result.Error != nil {
		return result.Error
	}

	return nil
}

// Function to find a user in the database
func (us *UserService) Find(user models.User, id uint) error {
	result := us.DB.Preload("Posts").First(&user, id)

	// If an error occurs while searching for the user, return the error
	if result.Error != nil {
		return result.Error
	}

	// Otherwise, return nil
	return nil
}

// Function to find all users in the database
func (us *UserService) FindAll() ([]models.User, error) {
	var users []models.User

	// Preload the user's posts and find all users
	result := us.DB.Preload("Posts").Find(&users)

	// If an error occurs while searching for the users, return the error
	if result.Error != nil {
		return nil, result.Error
	}

	// Otherwise, return the users
	return users, nil
}

// Function to find a user by ID in the database
func (us *UserService) FindByID(id uint) (models.User, error) {
	var user models.User

	// Preload the user's posts and find the user by ID
	result := us.DB.Preload("Posts").First(&user, id)

	// If an error occurs while searching for the user, return the error
	if result.Error != nil {
		return user, result.Error
	}

	// Otherwise, return the user
	return user, nil
}

// Function to update a user in the database
func (us *UserService) Update(email string, user models.UpdateUser) (models.User, error) {
	userToUpdate, err := us.FindByEmail(email)

	if err != nil {
		return models.User{}, err
	}

	result := us.DB.Model(&userToUpdate).Select("first_name", "last_name").Updates(models.User{FirstName: user.FirstName, LastName: user.LastName})

	if result.Error != nil {
		return models.User{}, result.Error
	}

	return userToUpdate, nil
}

// Function to find a user by email in the database
func (us *UserService) FindByEmail(email string) (models.User, error) {
	var user models.User
	result := us.DB.Where("email = ?", email).First(&user)

	if result.Error != nil {
		return models.User{}, result.Error
	}

	return user, nil
}
