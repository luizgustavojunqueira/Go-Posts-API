package main

import (
	"fmt"
	"github.com/joho/godotenv"
	"luizg/PostsAPI/api"
	"os"
)

func main() {

	err := godotenv.Load(".env")

	if err != nil {
		fmt.Println("Error loading .env file")
		return
	}

	api := &api.Api{}
	port := os.Getenv("API_PORT")
	api.Run(":" + port)
}
