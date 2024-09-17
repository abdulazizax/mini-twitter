package api

import (
	"context"
	"log"
	"log/slog"
	"net"
	"sync"

	"github.com/abdulazizax/mini-twitter/user-service/internal/items/msgbroker"
	"github.com/abdulazizax/mini-twitter/user-service/internal/items/service"
	"github.com/abdulazizax/mini-twitter/user-service/internal/items/storage"
	"github.com/abdulazizax/mini-twitter/user-service/internal/pkg/config"

	pb "github.com/abdulazizax/mini-twitter/user-service/genproto/user"

	sq "github.com/Masterminds/squirrel"
	"google.golang.org/grpc"
)

func Run(config *config.Config, logger *slog.Logger) error {
	db, err := storage.ConnectDB(config)
	if err != nil {
		logger.Error("Error while connecting to Postgres", slog.String("err", err.Error()))
		return err
	}

	sqrl := sq.StatementBuilder.PlaceholderFormat(sq.Dollar)

	service := service.New(storage.New(
		db,
		sqrl,
		config,
		logger,
	), logger)

	var wg sync.WaitGroup

	consumers := []struct {
		topic string
	}{
		{"user.registered"},
		{"user.updated"},
		{"user.deleted"},
		{"user.followed"},
		{"user.unfollowed"},
		{"user.password.updated"},
		{"user.verifycode.sended"},
	}

	for _, consumer := range consumers {
		c, err := msgbroker.NewConsumer(config.Kafka, "todolist", consumer.topic, logger, service, &wg)
		if err != nil {
			logger.Error("Error creating "+consumer.topic+" consumer", slog.String("err", err.Error()))
			continue
		}
		go c.Start(context.Background())
	}

	listener, err := net.Listen("tcp", config.Server.Port)
	if err != nil {
		logger.Error("Error while starting server", slog.String("err", err.Error()))
		return err
	}

	server := grpc.NewServer()
	pb.RegisterUserServiceServer(server, service)

	logger.Info("Server has started running", slog.String("port", config.Server.Port))
	log.Printf("Server has started running on port %s", config.Server.Port)

	return server.Serve(listener)
}
