package handlers

import (
	"luizg/PostsAPI/cmd/api/middlewares"
	"luizg/PostsAPI/internal/service"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type UserController struct {
	UserService *service.UserService
}

// Initialize User routes
func (controller *UserController) SetRoutes(router *gin.RouterGroup) {
	router.GET("/users", controller.GetUsers)
	router.GET("/users/:id", middlewares.AuthMiddleware(), controller.GetUser)
	router.GET("/users/me", middlewares.AuthMiddleware(), controller.GetCurrentUser)
	router.DELETE("/users", middlewares.AuthMiddleware(), controller.DeleteUser)
	router.PUT("/users", middlewares.AuthMiddleware(), controller.UpdateUser)
}

// Endpoint to get the current user
//
// @Summary Get the current user
// @Description Get the current user
// @Tags users
// @Produce json
// @Success 200 {object} service.User "Current user"
// @Failure 500 {string} string "Internal server error"
// @Router /users/current [get]
func (controller *UserController) GetCurrentUser(c *gin.Context) {
	userID := c.GetUint("user_id")

	if userID == 0 {
		panic("Something went wrong with the auth middleware")
	}

	user, err := controller.UserService.FindByID(userID)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not find user"})
		return
	}

	user.Password = ""

	c.JSON(http.StatusOK, gin.H{"user": user})
}

// Endpoint to get a user by ID
//
// @Summary Get a user by ID
// @Description Get a user by ID
// @Tags users
// @Produce json
// @Param id path int true "User ID"
// @Success 200 {object} service.User "User"
// @Failure 400 {string} string "Bad request"
// @Failure 500 {string} string "Internal server error"
// @Router /users/{id} [get]
func (controller *UserController) GetUser(c *gin.Context) {
	userIDInt, err := strconv.Atoi(c.Param("id"))

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
		return
	}

	userID := uint(userIDInt)

	user, err := controller.UserService.FindByID(userID)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not find user"})
		return
	}

	user.Password = ""

	c.JSON(http.StatusOK, gin.H{"user": user})
}

// Endpoint to get all users
//
// @Summary Get all users
// @Description Get all users
// @Tags users
// @Produce json
// @Success 200 {object} service.User "Users"
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
// @Param userUpdate body service.UpdateUser true "User update"
// @Success 200 {object} service.User "Updated user"
// @Failure 400 {string} string "Bad request"
// @Failure 500 {string} string "Internal server error"
// @Router /users [put]
func (controller *UserController) UpdateUser(c *gin.Context) {
	var updateUser service.UpdateUser

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
