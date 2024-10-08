package main

import (
	"log"
	"log/slog"
	"os"

	"github.com/abdulazizax/mini-twitter/user-service/cmd/api"
	"github.com/abdulazizax/mini-twitter/user-service/internal/pkg/config"
)

func main() {
	config, err := config.New()
	if err != nil {
		log.Fatalln(err)
	}

	logFile, err := os.OpenFile("application.log", os.O_TRUNC|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Fatal(err)
	}
	defer logFile.Close()

	logger := slog.New(slog.NewJSONHandler(logFile, nil))

	log.Fatalln(api.Run(config, logger))

}
