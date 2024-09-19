package userhandler

import (
	"log"

	"github.com/abdulazizax/mini-twitter/api-service/internal/items/msgbroker"
	"github.com/abdulazizax/mini-twitter/api-service/internal/pkg/config"
	"github.com/minio/minio-go/v7"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	pb "github.com/abdulazizax/mini-twitter/api-service/genproto/user"

	"log/slog"
)

type UserHandler struct {
	user     pb.UserServiceClient
	logger   *slog.Logger
	config   *config.Config
	producer *msgbroker.Producer
	minio    *minio.Client
}

func NewUserHandler(logger *slog.Logger, config *config.Config, producer *msgbroker.Producer, minio *minio.Client) *UserHandler {
	return &UserHandler{
		user:     pb.NewUserServiceClient(connect("localhost", config.Server.UserPort)),
		logger:   logger,
		config:   config,
		producer: producer,
		minio:    minio,
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
