package api

import (
	"luizg/PostsAPI/api/controllers"
	"luizg/PostsAPI/api/models"
	"luizg/PostsAPI/api/services"

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
		panic("failed to connect database")
	}

	api.DB = db

	api.DB.AutoMigrate(&models.User{})
	api.DB.AutoMigrate(&models.Post{})

	//Controllers
	userController := &controllers.UserController{UserService: &services.UserService{DB: api.DB}}
	userController.SetRoutes(api.Router)

	postController := &controllers.PostController{PostService: &services.PostService{DB: api.DB}}
	postController.SetRoutes(api.Router)
}

func (api *Api) Run(addr string) {
	api.Initialize()

	api.Router.Run(addr)
}
