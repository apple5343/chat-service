package integration_test

import (
	api "chat-service/pkg/api/chat_v1"
	"chat-service/test/integration/grpc"
	"context"
	"testing"
	"time"

	"github.com/brianvoe/gofakeit/v6"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const (
	testTimeout = 10 * time.Second
)

func TestAddUserHappyPath(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), testTimeout)
	defer cancel()
	client := grpc.NewClient(t)
	createChatRequest := api.CreateChatRequest{
		ProjectId: gofakeit.UUID(),
		Name:      gofakeit.Name(),
		Member:    []string{},
	}

	_, err := client.CreateChat(ctx, &createChatRequest)
	require.NoError(t, err)

	members := []string{}
	for i := 0; i < 10; i++ {
		id := gofakeit.UUID()
		members = append(members, id)
		r, err := client.AddUserToChat(ctx, &api.AddUserToChatRequest{
			ProjectId: createChatRequest.GetProjectId(),
			UserId:    id,
		})
		require.NoError(t, err)
		require.Equal(t, createChatRequest.GetProjectId(), r.GetProjectId())
	}

	resp, err := client .GetChat(ctx, &api.GetChatRequest{
		ProjectId: createChatRequest.ProjectId,
	})
	require.NoError(t, err)
	require.NotNil(t, resp)

	assert.Equal(t, createChatRequest.ProjectId, resp.Chat.GetProjectId())
	assert.Equal(t, createChatRequest.Name, resp.Chat.GetName())
	assert.Equal(t, members, resp.Chat.GetMembers())
}
