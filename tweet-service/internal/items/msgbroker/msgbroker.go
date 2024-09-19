package msgbroker

import (
	"context"
	"log/slog"
	"sync"

	"github.com/abdulazizax/mini-twitter/tweet-service/genproto/comment"
	"github.com/abdulazizax/mini-twitter/tweet-service/genproto/like"
	"github.com/abdulazizax/mini-twitter/tweet-service/genproto/tweet"
	"github.com/abdulazizax/mini-twitter/tweet-service/internal/items/service"
	"github.com/confluentinc/confluent-kafka-go/kafka"
	"google.golang.org/protobuf/encoding/protojson"
	"google.golang.org/protobuf/proto"
)

// Consumer structure
type Consumer struct {
	service  *service.Service
	consumer *kafka.Consumer
	logger   *slog.Logger
	wg       *sync.WaitGroup
}

// NewConsumer creates a new consumer
func NewConsumer(brokers, topic string, logger *slog.Logger, service *service.Service, wg *sync.WaitGroup) (*Consumer, error) {
	consumer, err := kafka.NewConsumer(&kafka.ConfigMap{
		"bootstrap.servers": brokers,
		"auto.offset.reset": "earliest",
		"group.id":          "tweet-service-group",
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

// Start starts the consumer to process messages
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
			case "user.tweet.created":
				var req tweet.CreateTweetRequest
				errUnmarshal = protojson.Unmarshal(msg.Value, &req)
				response, err = c.service.TweetService.CreateTweet(ctx, &req)
			case "user.tweet.updated":
				var req tweet.UpdateTweetRequest
				errUnmarshal = protojson.Unmarshal(msg.Value, &req)
				response, err = c.service.TweetService.UpdateTweet(ctx, &req)
			case "user.tweet.deleted":
				var req tweet.DeleteTweetRequest
				errUnmarshal = protojson.Unmarshal(msg.Value, &req)
				response, err = c.service.TweetService.DeleteTweet(ctx, &req)
			case "user.tweet.views.increased":
				var req tweet.Id
				errUnmarshal = protojson.Unmarshal(msg.Value, &req)
				response, err = c.service.TweetService.IncreaseViewsCount(ctx, &req)
			case "user.tweet.repost.increased":
				var req tweet.Id
				errUnmarshal = protojson.Unmarshal(msg.Value, &req)
				response, err = c.service.TweetService.IncreaseRepostCount(ctx, &req)
			case "user.tweet.shares.increased":
				var req tweet.Id
				errUnmarshal = protojson.Unmarshal(msg.Value, &req)
				response, err = c.service.TweetService.IncreaseSharesCount(ctx, &req)
			case "user.tweet.comment.created":
				var req comment.CreateCommentRequest
				errUnmarshal = protojson.Unmarshal(msg.Value, &req)
				response, err = c.service.CommentService.CreateComment(ctx, &req)
			case "user.tweet.comment.deleted":
				var req comment.DeleteCommentRequest
				errUnmarshal = protojson.Unmarshal(msg.Value, &req)
				response, err = c.service.CommentService.DeleteComment(ctx, &req)
			case "user.tweet.like":
				var req like.LikeRequest
				errUnmarshal = protojson.Unmarshal(msg.Value, &req)
				response, err = c.service.LikeService.Like(ctx, &req)
			case "user.tweet.unlike":
				var req like.UnlikeRequest
				errUnmarshal = protojson.Unmarshal(msg.Value, &req)
				response, err = c.service.LikeService.Unlike(ctx, &req)
			case "user.tweet.comment.like":
				var req like.LikeRequest
				errUnmarshal = protojson.Unmarshal(msg.Value, &req)
				response, err = c.service.LikeService.Like(ctx, &req)
			case "user.tweet.comment.unlike":
				var req like.UnlikeRequest
				errUnmarshal = protojson.Unmarshal(msg.Value, &req)
				response, err = c.service.LikeService.Unlike(ctx, &req)
			default:
				c.logger.Warn("Unknown topic", "topic", *msg.TopicPartition.Topic)
			}

			if errUnmarshal != nil {
				c.logger.Error("Error while unmarshaling data", "error", errUnmarshal)
				continue
			}

			if err != nil {
				c.logger.Error("Failed processing message", "topic", *msg.TopicPartition.Topic, "error", err.Error())
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
