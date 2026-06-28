package integration_test

import (
	"context"
	"testing"

	api "chat-service/pkg/api/chat_v1"
	"chat-service/test/integration/grpc"

	"github.com/brianvoe/gofakeit/v6"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestCreateChatHappyPath(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), testTimeout)
	defer cancel()
	client := grpc.NewClient(t)
	req := api.CreateChatRequest{
		ProjectId: gofakeit.UUID(),
		Name:      gofakeit.Name(),
		Member:    []string{gofakeit.UUID()},
	}

	_, err := client.CreateChat(ctx, &req)
	require.NoError(t, err)

	resp, err := client.GetChat(ctx, &api.GetChatRequest{
		ProjectId: req.ProjectId,
	})
	require.NoError(t, err)
	require.NotNil(t, resp)

	assert.Equal(t, req.ProjectId, resp.Chat.ProjectId)
	assert.Equal(t, req.Name, resp.Chat.Name)
	assert.Equal(t, req.Member, resp.Chat.Members)
}

func TestCreateChatBadRequest(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), testTimeout)
	defer cancel()
	client := grpc.NewClient(t)
	_, err := client.CreateChat(ctx, &api.CreateChatRequest{})
	require.Error(t, err)
	assert.Contains(t, err.Error(), "rpc error: code = InvalidArgument desc = bad request")
}
