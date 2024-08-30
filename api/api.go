package api

import (
	"luizg/PostsAPI/api/controllers"
	"luizg/PostsAPI/api/models"
	"luizg/PostsAPI/docs"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/swaggo/files"       // swagger embed files
	"github.com/swaggo/gin-swagger" // gin-swagger middleware
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

	config := cors.DefaultConfig()
	config.AllowOrigins = []string{"http://localhost:5173"}
	config.AllowCredentials = true
	config.AllowMethods = []string{"GET", "POST", "DELETE", "PUT", "OPTIONS"}
	config.AllowHeaders = []string{"Origin", "Content-Length", "Content-Type", "Authorization"}
	api.Router.Use(cors.New(config))

	v1 := api.Router.Group("/api/v1")

	//Controllers
	userController := &controllers.UserController{UserService: &models.UserService{DB: api.DB}}
	userController.SetRoutes(v1)

	postController := &controllers.PostController{PostService: &models.PostService{DB: api.DB}}
	postController.SetRoutes(v1)

	authController := &controllers.AuthController{UserService: &models.UserService{DB: api.DB}}
	authController.SetRoutes(v1)

	api.Router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
}

func (api *Api) Run(addr string) {
	api.Initialize()

	docs.SwaggerInfo.Title = "Posts API"
	docs.SwaggerInfo.Description = "A simple API to create posts and users"
	docs.SwaggerInfo.Version = "1.0"
	docs.SwaggerInfo.Host = "localhost" + addr
	docs.SwaggerInfo.BasePath = "/api/v1"
	docs.SwaggerInfo.Schemes = []string{"http"}

	api.Router.Run(addr)
}
