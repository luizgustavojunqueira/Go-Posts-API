package main

import (
	"luizg/PostsAPI/api"
)

func main() {
	api := &api.Api{}
	api.Run(":8080")
}
