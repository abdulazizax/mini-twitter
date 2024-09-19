package tweethandler

import (
	"log"
	"log/slog"

	"github.com/abdulazizax/mini-twitter/api-service/genproto/comment"
	"github.com/abdulazizax/mini-twitter/api-service/genproto/like"
	"github.com/abdulazizax/mini-twitter/api-service/genproto/tweet"

	"github.com/abdulazizax/mini-twitter/api-service/internal/items/msgbroker"
	"github.com/abdulazizax/mini-twitter/api-service/internal/pkg/config"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type TwetterClientConn struct {
	TweetClient   tweet.TweetServiceClient
	CommentClient comment.CommentServiceClient
	LikeClient    like.LikeServiceClient
}

func NewTwetterClientConn(config *config.Config) *TwetterClientConn {
	return &TwetterClientConn{
		TweetClient:   tweet.NewTweetServiceClient(connect("localhost", config.Server.TwitterPort)),
		CommentClient: comment.NewCommentServiceClient(connect("localhost", config.Server.TwitterPort)),
		LikeClient:    like.NewLikeServiceClient(connect("localhost", config.Server.TwitterPort)),
	}
}

type TwitterHandler struct {
	CommentHandler *CommentHandler
	LikeHandler    *LikeHandler
	TweetHandler   *TweetHandler
}

func NewTwitterHandler(logger *slog.Logger, config *config.Config, producer *msgbroker.Producer) *TwitterHandler {
	clientConn := NewTwetterClientConn(config)

	return &TwitterHandler{
		CommentHandler: NewCommentHandler(clientConn.CommentClient, logger, config, producer),
		LikeHandler:    NewLikeHandler(clientConn.LikeClient, logger, config, producer),
		TweetHandler:   NewTweetHandler(clientConn.TweetClient, logger, config, producer),
	}
}

func connect(host, port string) *grpc.ClientConn {
	address := host + port
	conn, err := grpc.NewClient(address, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatal(err)
	}
	return conn
}
