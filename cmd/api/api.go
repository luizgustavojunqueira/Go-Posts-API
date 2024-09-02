package api

import (
	"luizg/PostsAPI/cmd/api/handlers"
	"luizg/PostsAPI/docs"
	"luizg/PostsAPI/internal/service"

	"os"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"     // swagger embed files
	ginSwagger "github.com/swaggo/gin-swagger" // gin-swagger middleware
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
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

	errMigrate := api.DB.AutoMigrate(&service.User{}, &service.Post{})

	if errMigrate != nil {
		panic("Failed to migrate database")
	}

	config := cors.DefaultConfig()
	config.AllowOrigins = []string{"http://localhost:5173"}
	config.AllowCredentials = true
	config.AllowMethods = []string{"GET", "POST", "DELETE", "PUT", "OPTIONS"}
	config.AllowHeaders = []string{"Origin", "Content-Length", "Content-Type", "Authorization"}
	api.Router.Use(cors.New(config))

	v1 := api.Router.Group("/api/v1")

	//Controllers
	userController := &handlers.UserController{UserService: &service.UserService{DB: api.DB}}
	userController.SetRoutes(v1)

	postController := &handlers.PostController{PostService: &service.PostService{DB: api.DB}}
	postController.SetRoutes(v1)

	authController := &handlers.AuthController{UserService: &service.UserService{DB: api.DB}}
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

	err := api.Router.Run(addr)

	if err != nil {
		panic("Failed to run server")
	}
}
