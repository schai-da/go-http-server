package main

import (
	"go-http-server/internal/server"
	"log"
)

const (
	address = "localhost:12345"
)

func main() {
	server := server.HttpServer{
		Address: address,
	}

	err := server.Start()
	if err != nil {
		log.Fatalf("Failed to start HTTP server: %v", err)
	}
}
