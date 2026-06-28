package redis

import (
	"chat-service/entity"
	"context"
	"fmt"
	"testing"

	"github.com/redis/go-redis/v9"
)

type RedisClient struct {
	redis  *redis.Client
	config Config
}

type Config struct {
	Host       string `env:"REDIS_HOST" env-required:"true"`
	Port       string `env:"REDIS_PORT" env-required:"true"`
	StreamName string `env:"STREAM_NAME" env-required:"true"`
}

func NewClient(t *testing.T, cfg Config) *RedisClient {
	t.Helper()
	client := redis.NewClient(&redis.Options{
		Addr: fmt.Sprintf("%s:%s", cfg.Host, cfg.Port),
	})
	if err := client.Ping(context.Background()).Err(); err != nil {
		t.Fatal(err)
	}

	return &RedisClient{
		redis:  client,
		config: cfg,
	}
}

func (s *RedisClient) SendMessage(ctx context.Context, msg *entity.Message) error {
	return s.redis.XAdd(ctx, &redis.XAddArgs{
		Stream: s.config.StreamName,
		Values: map[string]string{
			"content": msg.Content,
			"user_id": msg.UserID,
			"room_id": msg.ProjectID,
		},
	}).Err()
}

func (s *RedisClient) Close() {
	s.redis.Close()
}
