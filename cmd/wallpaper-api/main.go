package main

import (
	"log"

	"github.com/joho/godotenv"
	server "github.com/wallpaper-api/internal"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("Error loading .env file: %s", err)
	}

	server := server.HttpServer{}

	server.StartServer()
}
