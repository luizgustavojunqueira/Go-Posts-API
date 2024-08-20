package controllers

import (
	"github.com/gin-gonic/gin"
	"luizg/PostsAPI/api/models"
	"luizg/PostsAPI/api/services"
	"luizg/PostsAPI/utils"
	"net/http"
)

type UserController struct {
	UserService *services.UserService
}

// Initialize User routes
func (controller *UserController) SetRoutes(router *gin.Engine) {
	router.POST("/users", controller.CreateUser)
	router.GET("/users", controller.GetUsers)
}

// Endpoint to create a new user
func (controller *UserController) CreateUser(c *gin.Context) {
	var user models.User

	// Parse the request body to the `user` struct
	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	}

	// Hash the user password
	passHash, err := utils.HashPassword(user.Password)

	// If an error occurs, return the error
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
	}

	user.Password = passHash

	// Save the user in the database
	id, err := controller.UserService.Save(user)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
	}

	// Respond with the new user ID
	c.JSON(http.StatusCreated, gin.H{"id": id})

}

// Endpoint to get all users
func (controller *UserController) GetUsers(c *gin.Context) {

	// Find all users
	users, err := controller.UserService.FindAll()

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
	}

	// Remove the password from the response
	for i := range users {
		users[i].Password = ""
	}

	// Respond with the users
	c.JSON(http.StatusOK, gin.H{"users": users})
}
