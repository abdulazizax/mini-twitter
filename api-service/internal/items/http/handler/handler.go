package handler

import (
	"log/slog"

	tweethandler "github.com/abdulazizax/mini-twitter/api-service/internal/items/http/handler/tweet-handler"
	userhandler "github.com/abdulazizax/mini-twitter/api-service/internal/items/http/handler/user-handler"
	"github.com/abdulazizax/mini-twitter/api-service/internal/pkg/config"
	"github.com/minio/minio-go/v7"

	"github.com/abdulazizax/mini-twitter/api-service/internal/items/msgbroker"
)

type Handler struct {
	UserHandler    *userhandler.UserHandler
	TwitterHandler *tweethandler.TwitterHandler
}

func New(logger *slog.Logger, config *config.Config, producer *msgbroker.Producer, minio *minio.Client) *Handler {

	return &Handler{
		UserHandler:    userhandler.NewUserHandler(logger, config, producer, minio),
		TwitterHandler: tweethandler.NewTwitterHandler(logger, config, producer),
	}
}
