package main

import (
	"os"

	"github.com/tg2648/trigo/internal/trivia"
)

func main() {
	staticDir := "./static"
	port := os.Getenv("TRIGO_PORT")
	if port == "" {
		port = "8080"
	}

	hub := trivia.NewHub()
	go hub.Run()

	server := trivia.NewServer(staticDir, hub)
	server.Run(port)
}
