package grpc

import (
	api "chat-service/pkg/api/chat_v1"
	"context"
	"fmt"
	"os"
	"testing"

	"github.com/joho/godotenv"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func NewClient(t *testing.T) api.ChatServiceClient {
	t.Helper()
	var grpcHost, grpcPort string
	if err := godotenv.Load("../../test.env"); err != nil {
		t.Fatal(err)
	}
	grpcHost = os.Getenv("GRPC_HOST")
	grpcPort = os.Getenv("GRPC_PORT")
	if grpcHost == "" || grpcPort == "" {
		t.Fatal("GRPC_HOST and GRPC_PORT must be set")
	}

	cc, err := grpc.DialContext(context.Background(), fmt.Sprintf("%s:%s", grpcHost, grpcPort), grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		t.Fatalf("grpc server connection failed: %v", err)
	}
	return api.NewChatServiceClient(cc)
}
