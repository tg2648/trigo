package main

import (
	"os"

	"github.com/tg2648/trigo/internal/trivia"
)

func main() {
	port := os.Getenv("TRIGO_PORT")
	if port == "" {
		port = "8080"
	}

	staticFolderDir := "./static"
	server := trivia.NewServer(staticFolderDir)
	server.Run(port)
}
