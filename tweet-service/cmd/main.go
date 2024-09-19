package main

import (
	"log"
	"log/slog"
	"os"

	"github.com/abdulazizax/mini-twitter/tweet-service/cmd/api"
	"github.com/abdulazizax/mini-twitter/tweet-service/internal/pkg/config"
)

func main() {
	// Load configuration
	config, err := config.New()
	if err != nil {
		log.Fatalln(err)
	}

	// Open log file
	logFile, err := os.OpenFile("application.log", os.O_TRUNC|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Fatal(err)
	}
	defer logFile.Close()

	// Create a new JSON logger
	logger := slog.New(slog.NewJSONHandler(logFile, nil))

	// Run the API and log any fatal errors
	log.Fatalln(api.Run(config, logger))
}
