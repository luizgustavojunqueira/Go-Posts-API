package controllers

import (
	"github.com/gin-gonic/gin"
	"luizg/PostsAPI/api/models"
	"luizg/PostsAPI/utils"
	"net/http"
	"net/mail"
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

	var newUser models.RegisterUser
	var errorMessages []string

	if err := c.ShouldBindJSON(&newUser); err != nil {
		errorMessages = append(errorMessages, err.Error())
	}

	if newUser.Password != newUser.ConfirmPassword {
		errorMessages = append(errorMessages, "Passwords do not match")
	}

	if len(newUser.Password) < 8 {
		errorMessages = append(errorMessages, "Password must have at least 8 characters")
	}

	if len(newUser.FirstName) < 2 {
		errorMessages = append(errorMessages, "First name must have at least 2 characters")
	}

	if len(newUser.LastName) < 2 {
		errorMessages = append(errorMessages, "Last name must have at least 2 characters")
	}

	if _, err := mail.ParseAddress(newUser.Email); err != nil {
		errorMessages = append(errorMessages, "Invalid email")
	}

	if len(errorMessages) > 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": errorMessages})
		return
	}

	passHash, err := utils.HashPassword(newUser.Password)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not hash password"})
		return
	}

	newUser.Password = passHash

	user := models.User{
		FirstName: newUser.FirstName,
		LastName:  newUser.LastName,
		Email:     newUser.Email,
		Password:  newUser.Password,
	}

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
