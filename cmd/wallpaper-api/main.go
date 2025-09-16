package main

import (
	server "github.com/wallpaper-api/internal"
)

func main() {
	server := server.HttpServer{}

	server.StartServer()
}
