package redisservice

import (
	"context"
	"fmt"
	"log/slog"
	"time"

	"github.com/abdulazizax/mini-twitter/user-service/internal/pkg/config"

	"github.com/go-redis/redis/v8"
)

type RedisService struct {
	redisDb *redis.Client
	logger  *slog.Logger
}

func New(redisDb *redis.Client, logger *slog.Logger) *RedisService {
	return &RedisService{
		logger:  logger,
		redisDb: redisDb,
	}
}

func NewRedisClient(cfg *config.Config) *redis.Client {
	rdb := redis.NewClient(&redis.Options{
		Addr:     cfg.RedisURI,
		Password: "",
		DB:       0,
	})

	return rdb
}

func (r *RedisService) StoreEmailAndCode(ctx context.Context, email string, code int) error {
	codeKey := "verification_code:" + email
	err := r.redisDb.Set(ctx, codeKey, code, time.Minute*15).Err()
	if err != nil {
		r.logger.Error("Error while storing verification code", slog.String("error", err.Error()))
		return err
	}
	return nil
}

func (r *RedisService) GetCodeByEmail(ctx context.Context, email string) (int, error) {
	codeKey := "verification_code:" + email
	codeStr, err := r.redisDb.Get(ctx, codeKey).Result()
	if err == redis.Nil {
		return 0, nil
	} else if err != nil {
		r.logger.Error("Error while getting verification code", slog.String("error", err.Error()))
		return 0, err
	}

	var code int
	_, err = fmt.Sscanf(codeStr, "%d", &code)
	if err != nil {
		r.logger.Error("Error while parsing verification code", slog.String("error", err.Error()))
		return 0, err
	}

	return code, nil
}
