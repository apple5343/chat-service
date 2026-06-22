package redis

import (
	"chat-service/entity"
	"context"
	"fmt"
	"os"
	"testing"
	"time"

	"github.com/redis/go-redis/v9"
)

const timeout = 5

type Suite struct {
	redis *redis.Client
	*testing.T
	streamName string
}

func New(t *testing.T) (context.Context, *Suite) {
	var redisHost, redisPort, streamName string
	t.Helper()
	redisHost, _ = os.LookupEnv("REDIS_HOST")
	redisPort, _ = os.LookupEnv("REDIS_PORT")
	streamName, _ = os.LookupEnv("STREAM_NAME")
	if redisHost == "" || redisPort == "" || streamName == "" {
		t.Fatal("REDIS_HOST, REDIS_PORT and STREAM_NAME must be set")
	}
	client := redis.NewClient(&redis.Options{
		Addr: fmt.Sprintf("%s:%s", redisHost, redisPort),
	})
	if err := client.Ping(context.Background()).Err(); err != nil {
		t.Fatal(err)
	}
	ctx, cancel := context.WithTimeout(context.Background(), timeout*time.Second)
	t.Cleanup(func() {
		t.Helper()
		cancel()
	})
	return ctx, &Suite{
		redis:      client,
		T:          t,
		streamName: streamName,
	}
}

func (s *Suite) SendMessage(ctx context.Context, msg *entity.Message) error {
	return s.redis.XAdd(ctx, &redis.XAddArgs{
		Stream: s.streamName,
		Values: map[string]string{
			"content": msg.Content,
			"user_id": msg.UserID,
			"room_id": msg.ProjectID,
		},
	}).Err()
}
