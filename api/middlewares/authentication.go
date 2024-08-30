package middlewares

import (
	"github.com/gin-gonic/gin"
	"luizg/PostsAPI/utils"
	"net/http"
	"strings"
)

// Middleware to authenticate the user
func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		var request struct {
			Authorization string `header:"Authorization"`
		}

		// Get token from header
		if err := c.ShouldBindHeader(&request); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Missing authorization header"})
			c.Abort()
			return
		}

		// Remove Bearer from token
		request.Authorization = strings.TrimPrefix(request.Authorization, "Bearer ")

		userClaims, err := utils.ParseToken(request.Authorization)

		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid Token"})
			c.Abort()
			return
		}

		// Set user id in context
		c.Set("user_id", userClaims.UserID)
		c.Next()
	}
}
