package integration_test

import (
	"chat-service/entity"
	api "chat-service/pkg/api/chat_v1"
	"chat-service/test/integration/grpc"
	"chat-service/test/integration/redis"
	"context"
	"testing"
	"time"

	"github.com/brianvoe/gofakeit/v6"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

const (
	readMessageTimeout = 2 * time.Second
	sendMessageTimeout = 200 * time.Millisecond
)

func TestSendMessagesHappyPath(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), testTimeout)
	defer cancel()
	grpcClient := grpc.NewClient(t)
	redisClient := redis.NewClient(t)

	createChatRequest := api.CreateChatRequest{
		ProjectId: gofakeit.UUID(),
		Name:      gofakeit.Name(),
		Member:    []string{},
	}

	_, err := grpcClient.CreateChat(ctx, &createChatRequest)
	require.NoError(t, err)

	userId := gofakeit.UUID()
	_, err = grpcClient.AddUserToChat(ctx, &api.AddUserToChatRequest{
		ProjectId: createChatRequest.GetProjectId(),
		UserId:    userId,
	})
	require.NoError(t, err)

	messages := []*entity.Message{}
	for i := 0; i < 10; i++ {
		msg := &entity.Message{
			ProjectID: createChatRequest.GetProjectId(),
			UserID:    userId,
			Content:   gofakeit.Word(),
		}
		messages = append(messages, msg)
		err := redisClient.SendMessage(ctx, msg)
		time.Sleep(sendMessageTimeout)
		require.NoError(t, err)
	}
	time.Sleep(readMessageTimeout)
	resp, err := grpcClient.GetMessages(ctx, &api.GetMessagesRequest{
		UserId:    userId,
		ProjectId: createChatRequest.GetProjectId(),
		Limit:     10,
		Cursor:    1,
	})
	require.NoError(t, err)
	require.NotNil(t, resp)
	require.Len(t, resp.GetMessages(), len(messages))
	for i := 0; i < 10; i++ {
		assert.Equal(t, messages[i].ProjectID, resp.Messages[len(messages)-1-i].ProjectId)
		assert.Equal(t, messages[i].UserID, resp.Messages[len(messages)-1-i].UserId)
		assert.Equal(t, messages[i].Content, resp.Messages[len(messages)-1-i].Content)
	}
}

func TestSendMessagesIncorrectUser(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), testTimeout)
	defer cancel()
	grpcClient := grpc.NewClient(t)
	redisClient := redis.NewClient(t)

	createChatRequest := api.CreateChatRequest{
		ProjectId: gofakeit.UUID(),
		Name:      gofakeit.Name(),
		Member:    []string{},
	}

	_, err := grpcClient.CreateChat(ctx, &createChatRequest)
	require.NoError(t, err)
	userIdFake := gofakeit.UUID()
	userId := gofakeit.UUID()

	_, err = grpcClient.AddUserToChat(ctx, &api.AddUserToChatRequest{
		ProjectId: createChatRequest.GetProjectId(),
		UserId:    userId,
	})
	require.NoError(t, err)
	//
	msg := &entity.Message{
		ProjectID: createChatRequest.GetProjectId(),
		UserID:    userId,
		Content:   gofakeit.Word(),
	}
	err = redisClient.SendMessage(ctx, msg)
	require.NoError(t, err)
	time.Sleep(readMessageTimeout)
	msgFake := &entity.Message{
		ProjectID: createChatRequest.GetProjectId(),
		UserID:    userIdFake,
		Content:   gofakeit.Word(),
	}
	err = redisClient.SendMessage(ctx, msgFake)
	require.NoError(t, err)
	resp, err := grpcClient.GetMessages(ctx, &api.GetMessagesRequest{
		UserId:    userId,
		ProjectId: createChatRequest.GetProjectId(),
		Limit:     10,
		Cursor:    1,
	})
	require.NoError(t, err)
	require.Equal(t, 1, len(resp.Messages))
	assert.Equal(t, msg.ProjectID, resp.Messages[0].ProjectId)
	assert.Equal(t, msg.UserID, resp.Messages[0].UserId)
	assert.Equal(t, msg.Content, resp.Messages[0].Content)

	resp, err = grpcClient.GetMessages(ctx, &api.GetMessagesRequest{
		UserId:    userIdFake,
		ProjectId: createChatRequest.GetProjectId(),
		Limit:     10,
		Cursor:    1,
	})

	if s, ok := status.FromError(err); ok {
		assert.EqualValues(t, codes.PermissionDenied, s.Code())
	}
}
