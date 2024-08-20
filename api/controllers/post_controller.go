package controllers

import (
	"github.com/gin-gonic/gin"
	"luizg/PostsAPI/api/models"
	"luizg/PostsAPI/api/services"
	"net/http"
	"strconv"
)

type PostController struct {
	PostService *services.PostService
}

func (controller *PostController) SetRoutes(router *gin.Engine) {
	router.POST("/posts", controller.CreatePost)
	router.GET("/posts", controller.GetPosts)
	router.DELETE("/posts/:id", controller.DeletePost)
}

func (controller *PostController) CreatePost(c *gin.Context) {
	var post models.Post

	if err := c.ShouldBindJSON(&post); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	id, err := controller.PostService.Save(post)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"id": id})
}

func (controller *PostController) GetPosts(c *gin.Context) {
	posts, err := controller.PostService.FindAll()

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"posts": posts})
}

func (controller *PostController) DeletePost(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	post, err := controller.PostService.FindByID(uint(id))

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	deletedId, err := controller.PostService.Delete(post)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"id": deletedId})

}
