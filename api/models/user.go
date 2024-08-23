package models

import (
	"gorm.io/gorm"
	"time"
)

type User struct {
	ID        uint      `gorm:"primaryKey"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	FirstName string    `json:"first_name"`
	LastName  string    `json:"last_name"`
	Email     string    `json:"email"`
	Password  string    `json:"password,omitempty"`
	Posts     []Post    `json:"posts" gorm:"constraint:OnDelete:CASCADE"`
}

type UserLogin struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type UpdateUser struct {
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
}

type UserService struct {
	DB *gorm.DB
}

// Function to save a user in the database
func (us *UserService) Save(user User) (uint, error) {
	result := us.DB.Create(&user)

	// If an error occurs while saving the user, return the error
	if result.Error != nil {
		return 0, result.Error
	}

	// Otherwise, return the new user ID
	return user.ID, nil
}

// Function to delete a user by email in the database
func (us *UserService) Delete(id uint) error {
	result := us.DB.Where("id = ?", id).Delete(&User{})

	if result.Error != nil {
		return result.Error
	}

	return nil
}

// Function to find a user in the database
func (us *UserService) Find(user User, id uint) error {
	result := us.DB.Preload("Posts").First(&user, id)

	// If an error occurs while searching for the user, return the error
	if result.Error != nil {
		return result.Error
	}

	// Otherwise, return nil
	return nil
}

// Function to find all users in the database
func (us *UserService) FindAll() ([]User, error) {
	var users []User

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
func (us *UserService) FindByID(id uint) (User, error) {
	var user User

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
func (us *UserService) Update(id uint, user UpdateUser) (User, error) {
	userToUpdate, err := us.FindByID(id)

	if err != nil {
		return User{}, err
	}

	result := us.DB.Model(&userToUpdate).Select("first_name", "last_name").Updates(User{FirstName: user.FirstName, LastName: user.LastName})

	if result.Error != nil {
		return User{}, result.Error
	}

	return userToUpdate, nil
}

// Function to find a user by email in the database
func (us *UserService) FindByEmail(email string) (User, error) {
	var user User
	result := us.DB.Where("email = ?", email).First(&user)

	if result.Error != nil {
		return User{}, result.Error
	}

	return user, nil
}
