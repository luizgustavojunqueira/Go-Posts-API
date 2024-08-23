package controllers

import (
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt"
	"luizg/PostsAPI/api/models"
	"luizg/PostsAPI/utils"
	"net/http"
	"time"
)

type AuthController struct {
	UserService *models.UserService
}

// Initialize Auth routes
func (controller *AuthController) SetRoutes(router *gin.Engine) {
	router.POST("/auth/login", controller.login)
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

	accessToken, err := utils.NewToken(&utils.UserClaims{
		UserID: user.ID,
		StandardClaims: jwt.StandardClaims{
			IssuedAt:  time.Now().Unix(),
			ExpiresAt: time.Now().Add(time.Minute * 15).Unix(),
		},
	})

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.Header("Authorization", "Bearer "+accessToken)
	c.JSON(http.StatusOK, gin.H{"message": "Logged in"})
}
