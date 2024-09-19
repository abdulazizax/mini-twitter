package main

import (
	"log"
	"log/slog"
	"os"
	"path/filepath"
	"time"

	casbin "github.com/casbin/casbin/v2"
	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"

	"github.com/abdulazizax/mini-twitter/api-service/internal/items/http/app"
	"github.com/abdulazizax/mini-twitter/api-service/internal/items/http/handler"
	"github.com/abdulazizax/mini-twitter/api-service/internal/items/msgbroker"
	"github.com/abdulazizax/mini-twitter/api-service/internal/pkg/config"
)

func main() {
	// Load configuration and handle errors
	config, err := config.New()
	if err != nil {
		log.Fatal(err)
	}

	// Set up logger
	logFile, err := os.OpenFile("application.log", os.O_TRUNC|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Fatal(err)
	}
	defer logFile.Close()

	logger := slog.New(slog.NewJSONHandler(logFile, nil))

	// Set up Casbin enforcer
	modelPath := filepath.Join("internal", "pkg", "casbin", "model.conf")
	policyPath := filepath.Join("internal", "pkg", "casbin", "policy.csv")

	enforcer, err := casbin.NewEnforcer(modelPath, policyPath)
	if err != nil {
		log.Fatal(err)
	}

	time.Sleep(2 * time.Second)

	log.Println(config.Kafka.Brokers)

	// Create Kafka producer
	producer, err := msgbroker.NewProducer(config, logger)
	if err != nil {
		log.Fatal(err)
	}
	defer producer.Close()

	// Create Kafka topics
	topics := []string{
		"user.registered",
		"user.updated",
		"user.deleted",
		"user.password.updated",

		"user.followed",
		"user.unfollowed",

		"user.tweet.created",
		"user.tweet.updated",
		"user.tweet.deleted",
		"user.tweet.views.increased",
		"user.tweet.repost.increased",
		"user.tweet.shares.increased",

		"user.tweet.comment.created",
		"user.tweet.comment.deleted",

		"user.tweet.like",
		"user.tweet.unlike",
		"user.tweet.comment.like",
		"user.tweet.comment.unlike",
	}
	err = producer.CreateTopics(topics, 1, 1) // Topic names, Number of partitions, Replication factor
	if err != nil {
		log.Fatal(err)
	}

	minio, err := minio.New(config.Minio, &minio.Options{
		Creds:  credentials.NewStaticV4("minioadmin", "minioadmin", ""),
		Secure: false,
	})
	if err != nil {
		log.Fatalln(err)
	}

	// Create handler
	handler := handler.New(logger, config, producer, minio)

	// Start HTTP server
	log.Fatal(app.Run(handler, logger, config, enforcer))
}
