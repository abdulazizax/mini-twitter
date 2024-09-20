package msgbroker

import (
	"context"
	"log/slog"
	"time"

	"github.com/IBM/sarama"
	"github.com/abdulazizax/mini-twitter/api-service/internal/pkg/config"
)

// Producer structure
type Producer struct {
	producer sarama.SyncProducer
	logger   *slog.Logger
	config   *config.Config
}

// NewProducer creates a new producer
func NewProducer(config *config.Config, logger *slog.Logger) (*Producer, error) {
	saramaConfig := sarama.NewConfig()
	saramaConfig.Producer.Return.Successes = true
	saramaConfig.Producer.RequiredAcks = sarama.WaitForAll // Ensure all in-sync replicas acknowledge the write

	producer, err := sarama.NewSyncProducer([]string{config.Kafka.Brokers}, saramaConfig)
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
	msg := &sarama.ProducerMessage{
		Topic: topic,
		Value: sarama.ByteEncoder(value),
	}

	// Produce the message
	partition, offset, err := p.producer.SendMessage(msg)
	if err != nil {
		p.logger.Error("Error producing message", "error", err.Error())
		return err
	}

	p.logger.Info("Message delivered", "topic", topic, "partition", partition, "offset", offset)
	return nil
}

// Close closes the producer
func (p *Producer) Close() {
	if err := p.producer.Close(); err != nil {
		p.logger.Error("Error closing producer", "error", err.Error())
	}
}

// CreateTopics creates multiple Kafka topics
func (p *Producer) CreateTopics(topics []string, numPartitions int, replicationFactor int) error {
	adminClient, err := sarama.NewClusterAdmin([]string{p.config.Kafka.Brokers}, sarama.NewConfig())
	if err != nil {
		return err
	}
	defer adminClient.Close()

	topicDetails := &sarama.TopicDetail{
		NumPartitions:     int32(numPartitions),
		ReplicationFactor: int16(replicationFactor),
	}

	_, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	for _, topic := range topics {
		// Check if topic exists
		topics, err := adminClient.ListTopics()
		if err != nil {
			p.logger.Error("Error listing topics", "error", err.Error())
			return err
		}
		if _, exists := topics[topic]; exists {
			p.logger.Warn("Topic already exists", "topic", topic)
			continue // If the topic exists, continue to the next one
		}

		// Create the topic
		if err := adminClient.CreateTopic(topic, topicDetails, false); err != nil {
			p.logger.Error("Error creating topic", "topic", topic, "error", err.Error())
			return err
		}
	}

	p.logger.Info("Successfully created topics", "topics", topics, "count", len(topics))
	return nil
}
