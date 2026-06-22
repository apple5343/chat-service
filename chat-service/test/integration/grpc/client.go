package grpc

import (
	api "chat-service/pkg/api/chat_v1"
	"context"
	"fmt"
	"os"
	"testing"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

const (
	timeout = 100
)

type Suite struct {
	*testing.T
	Client api.ChatServiceClient
}

func New(t *testing.T) (context.Context, *Suite) {
	var grpcHost, grpcPort string
	grpcHost, _ = os.LookupEnv("GRPC_HOST")
	grpcPort, _ = os.LookupEnv("GRPC_PORT")
	if grpcHost == "" || grpcPort == "" {
		t.Fatal("GRPC_HOST and GRPC_PORT must be set")
	}
	t.Helper()
	t.Parallel()
	ctx, cancel := context.WithTimeout(context.Background(), timeout*time.Second)
	t.Cleanup(func() {
		t.Helper()
		cancel()
	})

	cc, err := grpc.DialContext(context.Background(), fmt.Sprintf("%s:%s", grpcHost, grpcPort), grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		t.Fatalf("grpc server connection failed: %v", err)
	}
	return ctx, &Suite{t, api.NewChatServiceClient(cc)}
}
