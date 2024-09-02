package main

import (
	"fmt"
	"luizg/PostsAPI/cmd/api"
	"os"

	"github.com/joho/godotenv"
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
