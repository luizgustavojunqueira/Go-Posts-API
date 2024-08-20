package controllers

import (
	"github.com/gin-gonic/gin"
	"luizg/PostsAPI/api/models"
	"luizg/PostsAPI/api/services"
	"luizg/PostsAPI/utils"
	"net/http"
	"strconv"
)

type UserController struct {
	UserService *services.UserService
}

// Initialize User routes
func (controller *UserController) SetRoutes(router *gin.Engine) {
	router.POST("/users", controller.CreateUser)
	router.GET("/users", controller.GetUsers)
	router.DELETE("/users/:id", controller.DeleteUser)
	router.PUT("/users", controller.UpdateUser)
}

// Endpoint to create a new user
func (controller *UserController) CreateUser(c *gin.Context) {
	var user models.User

	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	passHash, err := utils.HashPassword(user.Password)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	user.Password = passHash

	id, err := controller.UserService.Save(user)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"id": id})

}

// Endpoint to get all users
func (controller *UserController) GetUsers(c *gin.Context) {

	users, err := controller.UserService.FindAll()

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	for i := range users {
		users[i].Password = ""
	}

	c.JSON(http.StatusOK, gin.H{"users": users})
}

// Encpoint to delete a user by id
func (controller *UserController) DeleteUser(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	user, err := controller.UserService.FindByID(uint(id))

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	deletedId, err := controller.UserService.Delete(user)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"id": deletedId})
}

// Endpoint to update a user
func (controller *UserController) UpdateUser(c *gin.Context) {
	var updatedUser models.User

	if err := c.ShouldBindJSON(&updatedUser); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	updatedUser, err := controller.UserService.Update(updatedUser)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"user": updatedUser})

}
