package controllers

import (
	"github.com/gin-gonic/gin"
	"luizg/PostsAPI/api/models"
	"luizg/PostsAPI/utils"
	"net/http"
)

type AuthController struct {
	UserService *models.UserService
}

// Initialize Auth routes
func (controller *AuthController) SetRoutes(router *gin.Engine) {
	router.POST("/auth/login", controller.login)
	router.POST("/auth/register", controller.register)
}

// Endpoint to login
func (controller *AuthController) login(c *gin.Context) {

	var userLogin models.UserLogin

	if err := c.ShouldBindJSON(&userLogin); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	user, err := controller.UserService.FindByEmail(userLogin.Email)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if !utils.CheckPasswordHash(userLogin.Password, user.Password) {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid password"})
		return
	}

	accessToken, err := utils.CreateTokenWithUserID(user.ID)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.Header("Authorization", "Bearer "+accessToken)
	c.JSON(http.StatusOK, gin.H{"message": "Logged in"})
}

// Endpoint to regirster a new user
func (controller *AuthController) register(c *gin.Context) {
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

	user_id, err := controller.UserService.Save(user)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not save user"})
		return
	}

	accessToken, err := utils.CreateTokenWithUserID(user_id)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.Header("Authorization", "Bearer "+accessToken)
	c.JSON(http.StatusCreated, gin.H{"message": "User registered"})
}
