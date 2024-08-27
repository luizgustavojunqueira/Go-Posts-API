package controllers

import (
	"github.com/gin-gonic/gin"
	"luizg/PostsAPI/api/middlewares"
	"luizg/PostsAPI/api/models"
	"net/http"
	"strconv"
)

type PostController struct {
	PostService *models.PostService
}

func (controller *PostController) SetRoutes(router *gin.Engine) {
	router.POST("/posts", middlewares.AuthMiddleware(), controller.CreatePost)
	router.GET("/posts", controller.GetPosts)
	router.DELETE("/posts/:id", middlewares.AuthMiddleware(), controller.DeletePost)
}

func (controller *PostController) CreatePost(c *gin.Context) {
	var postInput models.CreatePost

	if err := c.ShouldBindJSON(&postInput); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	userId := c.GetUint("user_id")

	if postInput.Title == "" || postInput.Content == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Title and Content are required"})
		return
	}

	if len(postInput.Title) > 100 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Title must be less than 100 characters"})
		return
	}

	if len(postInput.Content) > 500 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Content must be less than 500 characters"})
		return
	}

	newPost := models.Post{
		Title:   postInput.Title,
		Content: postInput.Content,
		UserID:  userId,
	}

	id, err := controller.PostService.Save(newPost)

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

	userId := c.GetUint("user_id")

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	post, err := controller.PostService.FindByID(uint(id))

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if post.UserID != userId {
		c.JSON(http.StatusForbidden, gin.H{"error": "You are not allowed to delete this post"})
		return
	}

	deletedId, err := controller.PostService.Delete(post)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"id": deletedId})

}
