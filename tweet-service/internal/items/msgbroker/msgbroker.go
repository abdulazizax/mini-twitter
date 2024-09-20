package msgbroker

import (
	"context"
	"log/slog"
	"sync"

	"github.com/IBM/sarama"
	"github.com/abdulazizax/mini-twitter/tweet-service/genproto/comment"
	"github.com/abdulazizax/mini-twitter/tweet-service/genproto/like"
	"github.com/abdulazizax/mini-twitter/tweet-service/genproto/tweet"
	"github.com/abdulazizax/mini-twitter/tweet-service/internal/items/service"
	"google.golang.org/protobuf/encoding/protojson"
	"google.golang.org/protobuf/proto"
)

// Consumer structure
type Consumer struct {
	service  *service.Service
	consumer sarama.Consumer
	logger   *slog.Logger
	wg       *sync.WaitGroup
}

// NewConsumer creates a new consumer
func NewConsumer(brokers, topic string, logger *slog.Logger, service *service.Service, wg *sync.WaitGroup) (*Consumer, error) {
	config := sarama.NewConfig()
	config.Consumer.Return.Errors = true

	consumer, err := sarama.NewConsumer([]string{brokers}, config)
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
func (c *Consumer) Start(ctx context.Context, topic string) {
	c.wg.Add(1)
	defer c.wg.Done()

	partitionList, err := c.consumer.Partitions(topic)
	if err != nil {
		c.logger.Error("Error getting partitions", "error", err.Error())
		return
	}

	var wg sync.WaitGroup

	for _, partition := range partitionList {
		pc, err := c.consumer.ConsumePartition(topic, partition, sarama.OffsetNewest)
		if err != nil {
			c.logger.Error("Error consuming partition", "partition", partition, "error", err.Error())
			continue
		}

		wg.Add(1)
		go func(pc sarama.PartitionConsumer) {
			defer wg.Done()
			defer pc.Close()

			for {
				select {
				case <-ctx.Done():
					c.logger.Info("Context done, stopping consumer")
					return
				case msg := <-pc.Messages():
					var response proto.Message
					var errUnmarshal error

					switch string(msg.Topic) {
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
						c.logger.Warn("Unknown topic", "topic", string(msg.Topic))
					}

					if errUnmarshal != nil {
						c.logger.Error("Error while unmarshaling data", "error", errUnmarshal)
						continue
					}

					if err != nil {
						c.logger.Error("Failed processing message", "topic", string(msg.Topic), "error", err.Error())
						continue
					}

					_, err = proto.Marshal(response)
					if err != nil {
						c.logger.Error("Failed to marshal response", "error", err)
						continue
					}

					c.logger.Info("Successfully processed message", "topic", string(msg.Topic))
				}
			}
		}(pc)
	}

	wg.Wait()
}
