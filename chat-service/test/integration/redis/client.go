package redis

import (
	"chat-service/entity"
	"context"
	"fmt"
	"os"
	"testing"

	"github.com/joho/godotenv"
	"github.com/redis/go-redis/v9"
)

type RedisClient struct {
	redis      *redis.Client
	streamName string
}

func NewClient(t *testing.T) *RedisClient {
	t.Helper()
	var redisHost, redisPort, streamName string
	if err := godotenv.Load("../../test.env"); err != nil {
		t.Fatal(err)
	}
	redisHost = os.Getenv("REDIS_HOST")
	redisPort = os.Getenv("REDIS_PORT")
	streamName = os.Getenv("STREAM_NAME")
	if redisHost == "" || redisPort == "" || streamName == "" {
		t.Fatal("REDIS_HOST, REDIS_PORT and STREAM_NAME must be set")
	}
	client := redis.NewClient(&redis.Options{
		Addr: fmt.Sprintf("%s:%s", redisHost, redisPort),
	})
	if err := client.Ping(context.Background()).Err(); err != nil {
		t.Fatal(err)
	}

	return &RedisClient{
		redis:      client,
		streamName: streamName,
	}
}

func (s *RedisClient) SendMessage(ctx context.Context, msg *entity.Message) error {
	return s.redis.XAdd(ctx, &redis.XAddArgs{
		Stream: s.streamName,
		Values: map[string]string{
			"content": msg.Content,
			"user_id": msg.UserID,
			"room_id": msg.ProjectID,
		},
	}).Err()
}
