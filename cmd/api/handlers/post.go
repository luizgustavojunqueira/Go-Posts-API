package handlers

import (
	"luizg/PostsAPI/cmd/api/middlewares"
	"luizg/PostsAPI/internal/service"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type PostController struct {
	PostService *service.PostService
}

func (controller *PostController) SetRoutes(router *gin.RouterGroup) {
	router.POST("/posts", middlewares.AuthMiddleware(), controller.CreatePost)
	router.GET("/posts/all", controller.GetPosts)
	router.GET("/posts/user", middlewares.AuthMiddleware(), controller.GetUserPosts)
	router.GET("/posts/user/:id", middlewares.AuthMiddleware(), controller.GetPostsByUserID)
	router.DELETE("/posts/:id", middlewares.AuthMiddleware(), controller.DeletePost)
}

// Endpoint to create a post
//
// @Summary Create a post
// @Description Create a post
// @Tags posts
// @Accept json
// @Produce json
// @Param post body service.CreatePost true "Post"
// @Success 201 {string} string "Created post"
// @Failure 400 {string} string "Bad request"
// @Failure 500 {string} string "Internal server error"
// @Router /posts [post]
func (controller *PostController) CreatePost(c *gin.Context) {
	var postInput service.CreatePost

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

	newPost := service.Post{
		Title:        postInput.Title,
		Content:      postInput.Content,
		UserID:       userId,
		UserFullName: c.GetString("user_full_name"),
	}

	id, err := controller.PostService.Save(newPost)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"id": id})
}

// Endpoint to get all posts
//
// @Summary Get all posts
// @Description Get all posts
// @Tags posts
// @Produce json
// @Success 200 {object} service.Post "Posts"
// @Failure 500 {string} string "Internal server error"
// @Router /posts/all [get]
func (controller *PostController) GetPosts(c *gin.Context) {
	posts, err := controller.PostService.FindAll()

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"posts": posts})
}

// Endpoint to get all posts from the logged user
//
// @Summary Get all posts from the logged user
// @Description Get all posts from the logged user
// @Tags posts
// @Produce json
// @Success 200 {object} service.Post "Posts"
// @Failure 500 {string} string "Internal server error"
// @Router /posts/user [get]
func (controller *PostController) GetUserPosts(c *gin.Context) {
	userId := c.GetUint("user_id")

	posts, err := controller.PostService.FindByUserID(userId)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"posts": posts})
}

// Endpoint to get all posts from a user
//
// @Summary Get all posts from a user
// @Description Get all posts from a user
// @Tags posts
// @Produce json
// @Param id path int true "User ID"
// @Success 200 {object} service.Post "Posts"
// @Failure 500 {string} string "Internal server error"
// @Router /posts/user/{id} [get]
func (controller *PostController) GetPostsByUserID(c *gin.Context) {
	userId, err := strconv.Atoi(c.Param("id"))

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	posts, err := controller.PostService.FindByUserID(uint(userId))

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"posts": posts})
}

// Endpoint to delete a post
//
// @Summary Delete a post
// @Description Delete a post
// @Tags posts
// @Produce json
// @Param id path int true "Post ID"
// @Success 200 {string} string "Deleted post"
// @Failure 500 {string} string "Internal server error"
// @Router /posts/{id} [delete]
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
