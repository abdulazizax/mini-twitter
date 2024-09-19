package api

import (
	"context"
	"log"
	"log/slog"
	"net"
	"sync"

	"github.com/abdulazizax/mini-twitter/tweet-service/internal/items/msgbroker"
	"github.com/abdulazizax/mini-twitter/tweet-service/internal/items/service"
	"github.com/abdulazizax/mini-twitter/tweet-service/internal/items/storage"
	psql "github.com/abdulazizax/mini-twitter/tweet-service/internal/items/storage/postgres"
	"github.com/abdulazizax/mini-twitter/tweet-service/internal/pkg/config"

	comment_pb "github.com/abdulazizax/mini-twitter/tweet-service/genproto/comment"
	like_pb "github.com/abdulazizax/mini-twitter/tweet-service/genproto/like"
	tweet_pb "github.com/abdulazizax/mini-twitter/tweet-service/genproto/tweet"

	sq "github.com/Masterminds/squirrel"
	"google.golang.org/grpc"
)

// Run initializes and starts the API server
func Run(config *config.Config, logger *slog.Logger) error {
	// Connect to the database
	db, err := psql.ConnectDB(config)
	if err != nil {
		logger.Error("Error while connecting to Postgres", slog.String("err", err.Error()))
		return err
	}

	// Initialize the SQL query builder
	sqrl := sq.StatementBuilder.PlaceholderFormat(sq.Dollar)

	// Create a new service instance
	service := service.New(storage.New(
		db,
		sqrl,
		config,
		logger,
	), logger)

	var wg sync.WaitGroup

	// Define Kafka consumer topics
	consumers := []struct {
		topic string
	}{
		{"user.tweet.created"},
		{"user.tweet.updated"},
		{"user.tweet.deleted"},
		{"user.tweet.views.increased"},
		{"user.tweet.repost.increased"},
		{"user.tweet.shares.increased"},
		{"user.tweet.comment.created"},
		{"user.tweet.comment.deleted"},
		{"user.tweet.like"},
		{"user.tweet.unlike"},
		{"user.tweet.comment.like"},
		{"user.tweet.comment.unlike"},
	}

	// Start Kafka consumers for each topic
	for _, consumer := range consumers {
		c, err := msgbroker.NewConsumer(config.Kafka, consumer.topic, logger, service, &wg)
		if err != nil {
			logger.Error("Error creating "+consumer.topic+" consumer", slog.String("err", err.Error()))
			continue
		}
		go c.Start(context.Background())
	}

	// Set up the gRPC server
	listener, err := net.Listen("tcp", config.Server.Port)
	if err != nil {
		logger.Error("Error while starting server", slog.String("err", err.Error()))
		return err
	}

	server := grpc.NewServer()
	tweet_pb.RegisterTweetServiceServer(server, service.TweetService)
	like_pb.RegisterLikeServiceServer(server, service.LikeService)
	comment_pb.RegisterCommentServiceServer(server, service.CommentService)

	// Log server start
	logger.Info("Server has started running", slog.String("port", config.Server.Port))
	log.Printf("Server has started running on port %s", config.Server.Port)

	// Start serving gRPC requests
	return server.Serve(listener)
}
