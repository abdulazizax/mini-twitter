package msgbroker

import (
	"context"
	"log/slog"
	"time"

	"github.com/abdulazizax/mini-twitter/api-service/internal/pkg/config"
	"github.com/confluentinc/confluent-kafka-go/kafka"
)

// Producer structure
type Producer struct {
	producer *kafka.Producer
	logger   *slog.Logger
	config   *config.Config
}

// NewProducer creates a new producer
func NewProducer(config *config.Config, logger *slog.Logger) (*Producer, error) {
	producer, err := kafka.NewProducer(&kafka.ConfigMap{
		"bootstrap.servers": config.Kafka.Brokers,
		"acks":              "all", // Ensure all in-sync replicas acknowledge the write
	})
	if err != nil {
		return nil, err
	}

	return &Producer{
		producer: producer,
		logger:   logger,
		config:   config,
	}, nil
}

// Send sends a message to a Kafka topic
func (p *Producer) Send(topic string, value []byte) error {
	msg := &kafka.Message{
		TopicPartition: kafka.TopicPartition{Topic: &topic, Partition: kafka.PartitionAny},
		Value:          value,
	}

	// Produce the message asynchronously
	err := p.producer.Produce(msg, nil)
	if err != nil {
		p.logger.Error("Error producing message", "error", err.Error())
		return err
	}

	// Wait for the delivery report
	ev := <-p.producer.Events()
	switch e := ev.(type) {
	case *kafka.Message:
		if e.TopicPartition.Error != nil {
			p.logger.Error("Failed to deliver message", "topic", *e.TopicPartition.Topic, "error", e.TopicPartition.Error.Error())
			return e.TopicPartition.Error
		}
		p.logger.Info("Message delivered", "topic", *e.TopicPartition.Topic, "partition", e.TopicPartition.Partition)
	default:
		p.logger.Warn("Received unsupported event type", "event", ev)
	}

	// Optionally, flush all messages after producing
	p.producer.Flush(15 * 1000) // 15 seconds
	return nil
}

// Close closes the producer
func (p *Producer) Close() {
	p.producer.Close()
}

// CreateTopics creates multiple Kafka topics
func (p *Producer) CreateTopics(topics []string, numPartitions int, replicationFactor int) error {
	adminClient, err := kafka.NewAdminClient(&kafka.ConfigMap{
		"bootstrap.servers": p.config.Kafka.Brokers,
	})
	if err != nil {
		return err
	}
	defer adminClient.Close()

	var topicConfigs []kafka.TopicSpecification
	for _, topic := range topics {
		topicConfigs = append(topicConfigs, kafka.TopicSpecification{
			Topic:             topic,
			NumPartitions:     numPartitions,
			ReplicationFactor: replicationFactor,
		})
	}

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	_, err = adminClient.CreateTopics(ctx, topicConfigs)
	if err != nil {
		p.logger.Error("Error creating topics", "error", err.Error())
		return err
	}

	p.logger.Info("Successfully created topics", "topics", topics, "count", len(topics))
	return nil
}
