package controllers

import (
	"github.com/gin-gonic/gin"
	"luizg/PostsAPI/api/middlewares"
	"luizg/PostsAPI/api/models"
	"net/http"
)

type UserController struct {
	UserService *models.UserService
}

// Initialize User routes
func (controller *UserController) SetRoutes(router *gin.RouterGroup) {
	router.GET("/users", controller.GetUsers)
	router.DELETE("/users", middlewares.AuthMiddleware(), controller.DeleteUser)
	router.PUT("/users", middlewares.AuthMiddleware(), controller.UpdateUser)
}

// Endpoint to get all users
//
// @Summary Get all users
// @Description Get all users
// @Tags users
// @Produce json
// @Success 200 {object} models.User "Users"
// @Failure 500 {string} string "Internal server error"
// @Router /users [get]
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

// Endpoint to delete a user
//
// @Summary Delete a user
// @Description Delete a user
// @Tags users
// @Produce json
// @Success 200 {string} string "Deleted user"
// @Failure 500 {string} string "Internal server error"
// @Router /users [delete]
func (controller *UserController) DeleteUser(c *gin.Context) {

	userID := c.GetUint("user_id")

	if userID == 0 {
		panic("Something went wrong with the auth middleware")
	}

	err := controller.UserService.Delete(userID)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not delete user"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"deleted_id": userID})
}

// Endpoint to update a user
//
// @Summary Update a user
// @Description Update a user
// @Tags users
// @Accept json
// @Produce json
// @Param userUpdate body models.UpdateUser true "User update"
// @Success 200 {object} models.User "Updated user"
// @Failure 400 {string} string "Bad request"
// @Failure 500 {string} string "Internal server error"
// @Router /users [put]
func (controller *UserController) UpdateUser(c *gin.Context) {
	var updateUser models.UpdateUser

	if err := c.ShouldBindJSON(&updateUser); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	userID := c.GetUint("user_id")

	if userID == 0 {
		panic("Something went wrong with the auth middleware")
	}

	updatedUser, err := controller.UserService.Update(userID, updateUser)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	updatedUser.Password = ""

	c.JSON(http.StatusOK, gin.H{"user": updatedUser})

}
