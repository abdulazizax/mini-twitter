package config

import (
	"os"

	"github.com/joho/godotenv"
)

type (
	Config struct {
		Server ServerConfig
		JWT    JWTConfig
		Kafka  KafkaConfig
		Minio  string
	}
	JWTConfig struct {
		SecretKey string
	}
	ServerConfig struct {
		ServerPort  string
		UserPort    string
		TwitterPort string
	}
	KafkaConfig struct {
		Brokers string
	}
)

func (c *Config) Load() error {
	if err := godotenv.Load(); err != nil {
		return err
	}

	c.Server.ServerPort = os.Getenv("SERVER_PORT")
	c.Server.UserPort = os.Getenv("USER_PORT")
	c.Server.TwitterPort = os.Getenv("TWITTER_PORT")
	c.JWT.SecretKey = os.Getenv("JWT_SECRET_KEY")
	c.Kafka.Brokers = os.Getenv("KAFKA_BROKER_URI")
	c.Minio = os.Getenv("MINIO_URI")

	return nil
}

func New() (*Config, error) {
	var config Config
	if err := config.Load(); err != nil {
		return nil, err
	}
	return &config, nil
}
