package grpc

import (
	api "chat-service/pkg/api/chat_v1"
	"context"
	"fmt"
	"testing"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type Config struct {
	Host string `env:"GRPC_HOST" env-required:"true"`
	Port string `env:"GRPC_PORT" env-required:"true"`
}

func NewClient(t *testing.T, cfg Config) api.ChatServiceClient {
	t.Helper()

	cc, err := grpc.DialContext(context.Background(), fmt.Sprintf("%s:%s", cfg.Host, cfg.Port), grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		t.Fatalf("grpc server connection failed: %v", err)
	}
	return api.NewChatServiceClient(cc)
}
