package controllers

import (
	"github.com/gin-gonic/gin"
	"luizg/PostsAPI/api/models"
	"luizg/PostsAPI/api/services"
	"net/http"
)

type PostController struct {
	PostService *services.PostService
}

func (controller *PostController) SetRoutes(router *gin.Engine) {
	router.POST("/posts", controller.CreatePost)
	router.GET("/posts", controller.GetPosts)
}

func (controller *PostController) CreatePost(c *gin.Context) {
	var post models.Post

	if err := c.ShouldBindJSON(&post); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	}

	id, err := controller.PostService.Save(post)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
	}

	c.JSON(http.StatusCreated, gin.H{"id": id})
}

func (controller *PostController) GetPosts(c *gin.Context) {
	posts, err := controller.PostService.FindAll()

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
	}

	c.JSON(http.StatusOK, gin.H{"posts": posts})
}
