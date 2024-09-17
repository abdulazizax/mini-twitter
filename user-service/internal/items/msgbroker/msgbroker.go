package msgbroker

import (
	"context"
	"log/slog"
	"sync"

	"github.com/abdulazizax/mini-twitter/user-service/internal/items/service"

	pb "github.com/abdulazizax/mini-twitter/user-service/genproto/user"

	"github.com/confluentinc/confluent-kafka-go/kafka"
	"google.golang.org/protobuf/encoding/protojson"
	"google.golang.org/protobuf/proto"
)

type Consumer struct {
	service  *service.Service
	consumer *kafka.Consumer
	logger   *slog.Logger
	wg       *sync.WaitGroup
}

func NewConsumer(brokers, groupID, topic string, logger *slog.Logger, service *service.Service, wg *sync.WaitGroup) (*Consumer, error) {
	consumer, err := kafka.NewConsumer(&kafka.ConfigMap{
		"bootstrap.servers": brokers,
		"group.id":          groupID,
		"auto.offset.reset": "earliest",
	})
	if err != nil {
		return nil, err
	}

	err = consumer.Subscribe(topic, nil)
	if err != nil {
		return nil, err
	}

	return &Consumer{
		service:  service,
		consumer: consumer,
		logger:   logger,
		wg:       wg,
	}, nil
}

func (c *Consumer) Start(ctx context.Context) {
	c.wg.Add(1)
	defer c.wg.Done()

	for {
		select {
		case <-ctx.Done():
			c.logger.Info("Context done, stopping consumer")
			c.consumer.Close()
			return
		default:
			msg, err := c.consumer.ReadMessage(-1)
			if err != nil {
				c.logger.Error("Error while consuming message", "error", err.Error())
				continue
			}

			var response proto.Message
			var errUnmarshal error

			switch *msg.TopicPartition.Topic {
			case "user.registered":
				var req pb.RegisterUserRequest
				errUnmarshal = protojson.Unmarshal(msg.Value, &req)
				response, err = c.service.RegisterUser(ctx, &req)
			case "user.updated":
				var req pb.UpdateUserRequest
				errUnmarshal = protojson.Unmarshal(msg.Value, &req)
				response, err = c.service.UpdateUser(ctx, &req)
			case "user.deleted":
				var req pb.DeleteUserRequest
				errUnmarshal = protojson.Unmarshal(msg.Value, &req)
				response, err = c.service.DeleteUser(ctx, &req)
			case "user.followed":
				var req pb.FollowUserRequest
				errUnmarshal = protojson.Unmarshal(msg.Value, &req)
				response, err = c.service.FollowUser(ctx, &req)
			case "user.unfollowed":
				var req pb.UnfollowUserRequest
				errUnmarshal = protojson.Unmarshal(msg.Value, &req)
				response, err = c.service.UnfollowUser(ctx, &req)
			case "user.password.updated":
				var req pb.UpdateUserPasswordRequest
				errUnmarshal = protojson.Unmarshal(msg.Value, &req)
				response, err = c.service.UpdateUserPassword(ctx, &req)
			case "user.verifycode.sended":
				var req pb.SendVerificationCodeRequest
				errUnmarshal = protojson.Unmarshal(msg.Value, &req)
				response, err = c.service.SendVerificationCode(ctx, &req)
			}

			if errUnmarshal != nil {
				c.logger.Error("Error while unmarshaling data", "error", errUnmarshal)
				continue
			}

			if err != nil {
				c.logger.Error("Failed in topic %s: %s", *msg.TopicPartition.Topic, err.Error())
				continue
			}

			_, err = proto.Marshal(response)
			if err != nil {
				c.logger.Error("Failed to marshal response", "error", err)
				continue
			}

			c.logger.Info("Successfully processed message", "topic", *msg.TopicPartition.Topic)
		}
	}
}
