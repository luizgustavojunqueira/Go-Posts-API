package api

import (
	"luizg/PostsAPI/api/controllers"
	"luizg/PostsAPI/api/models"

	"github.com/gin-gonic/gin"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"os"
)

type Api struct {
	Router *gin.Engine
	DB     *gorm.DB
}

func (api *Api) Initialize() {
	api.Router = gin.Default()

	db, err := gorm.Open(postgres.New(postgres.Config{
		DSN: os.Getenv("DATABASE_URL"),
	}), &gorm.Config{})

	if err != nil {
		panic("Failed to connect database")
	}

	api.DB = db

	api.DB.AutoMigrate(&models.User{}, &models.Post{})

	//Controllers
	userController := &controllers.UserController{UserService: &models.UserService{DB: api.DB}}
	userController.SetRoutes(api.Router)

	postController := &controllers.PostController{PostService: &models.PostService{DB: api.DB}}
	postController.SetRoutes(api.Router)

	authController := &controllers.AuthController{UserService: &models.UserService{DB: api.DB}}
	authController.SetRoutes(api.Router)
}

func (api *Api) Run(addr string) {
	api.Initialize()

	api.Router.Run(addr)
}
