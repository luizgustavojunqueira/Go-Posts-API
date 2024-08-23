package controllers

import (
	"github.com/gin-gonic/gin"
	"luizg/PostsAPI/api/middlewares"
	"luizg/PostsAPI/api/models"
	"luizg/PostsAPI/utils"
	"net/http"
)

type UserController struct {
	UserService *models.UserService
}

// Initialize User routes
func (controller *UserController) SetRoutes(router *gin.Engine) {
	router.POST("/users", controller.CreateUser)
	router.GET("/users", controller.GetUsers)
	router.DELETE("/users", middlewares.AuthMiddleware(), controller.DeleteUser)
	router.PUT("/users", middlewares.AuthMiddleware(), controller.UpdateUser)
}

// Endpoint to create a new user
func (controller *UserController) CreateUser(c *gin.Context) {
	var user models.User

	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
		return
	}

	passHash, err := utils.HashPassword(user.Password)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not hash password"})
		return
	}

	user.Password = passHash

	id, err := controller.UserService.Save(user)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not save user"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"id": id})

}

// Endpoint to get all users
func (controller *UserController) GetUsers(c *gin.Context) {

	users, err := controller.UserService.FindAll()

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not find users"})
		return
	}

	for i := range users {
		users[i].Password = ""
	}

	c.JSON(http.StatusOK, gin.H{"users": users})
}

// Encpoint to delete a user by id
func (controller *UserController) DeleteUser(c *gin.Context) {

	userEmail := c.GetString("user_email")

	if userEmail == "" {
		panic("Something went wrong with the auth middleware")
	}

	err := controller.UserService.Delete(userEmail)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not delete user"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"deleted_email": userEmail})
}

// Endpoint to update a user
func (controller *UserController) UpdateUser(c *gin.Context) {
	var updateUser models.UpdateUser

	if err := c.ShouldBindJSON(&updateUser); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	userEmail := c.GetString("user_email")

	if userEmail == "" {
		panic("Something went wrong with the auth middleware")
	}

	updatedUser, err := controller.UserService.Update(userEmail, updateUser)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	updatedUser.Password = ""

	c.JSON(http.StatusOK, gin.H{"user": updatedUser})

}
