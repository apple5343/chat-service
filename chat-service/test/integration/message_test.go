package integration

import (
	"chat-service/entity"
	api "chat-service/pkg/api/chat_v1"
	"chat-service/test/integration/grpc"
	"chat-service/test/integration/redis"
	"testing"
	"time"

	"github.com/brianvoe/gofakeit/v6"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
)

type MessageSuite struct {
	suite.Suite
	base        *BaseTestSuite
	grpcClient  api.ChatServiceClient
	redisClient *redis.RedisClient
}

func (s *MessageSuite) SetupSuite() {
	s.grpcClient = grpc.NewClient(s.T(), s.base.config.GrpcConfig)
	s.redisClient = redis.NewClient(s.T(), s.base.config.RedisConfig)
}

func (s *MessageSuite) TearDownSuite() {
	s.redisClient.Close()
}

func (s *MessageSuite) TestMessage() {
	s.T().Run("send messages", s.SendMessages)
}

func (s *MessageSuite) SendMessages(t *testing.T) {
	ctx := s.base.ctx

	createChatRequest := api.CreateChatRequest{
		ProjectId: gofakeit.UUID(),
		Name:      gofakeit.Name(),
		Member:    []string{},
	}

	_, err := s.grpcClient.CreateChat(ctx, &createChatRequest)
	require.NoError(t, err)

	userId := gofakeit.UUID()
	_, err = s.grpcClient.AddUserToChat(ctx, &api.AddUserToChatRequest{
		ProjectId: createChatRequest.GetProjectId(),
		UserId:    userId,
	})
	require.NoError(t, err)

	type compare struct {
		ProjectID string
		UserID    string
		Content   string
	}

	count := 10
	expected := make([]compare, 0, count)

	for i := 0; i < count; i++ {
		msg := &entity.Message{
			ProjectID: createChatRequest.GetProjectId(),
			UserID:    userId,
			Content:   gofakeit.Word(),
		}
		expected = append(expected, compare{
			ProjectID: msg.ProjectID,
			UserID:    msg.UserID,
			Content:   msg.Content,
		})
		err := s.redisClient.SendMessage(ctx, msg)
		require.NoError(t, err)
	}

	var resp *api.GetMessagesResponse
	length := 0
	require.Eventuallyf(t, func() bool {
		resp, err = s.grpcClient.GetMessages(ctx, &api.GetMessagesRequest{
			UserId:    userId,
			ProjectId: createChatRequest.GetProjectId(),
			Limit:     10,
			Cursor:    1,
		})
		if resp != nil {
			length = len(resp.Messages)
		}
		return err == nil && len(resp.Messages) == count
	}, 10*time.Second, 100*time.Millisecond, "expected %d messages, got %d", count, length)
	require.NoError(t, err)
	require.NotNil(t, resp)
	require.Len(t, resp.GetMessages(), len(expected))

	actual := make([]compare, 0, len(resp.GetMessages()))
	for i := 0; i < count; i++ {
		actual = append(actual, compare{
			ProjectID: resp.GetMessages()[i].GetProjectId(),
			UserID:    resp.GetMessages()[i].GetUserId(),
			Content:   resp.GetMessages()[i].GetContent(),
		})
	}

	require.ElementsMatch(t, expected, actual)
}

func (s *MessageSuite) SendMessagesInvalid(t *testing.T) {
	ctx := s.base.ctx

	createChatRequest := api.CreateChatRequest{
		ProjectId: gofakeit.UUID(),
		Name:      gofakeit.Name(),
		Member:    []string{},
	}

	_, err := s.grpcClient.CreateChat(ctx, &createChatRequest)
	require.NoError(t, err)
	userIdFake := gofakeit.UUID()
	userId := gofakeit.UUID()

	_, err = s.grpcClient.AddUserToChat(ctx, &api.AddUserToChatRequest{
		ProjectId: createChatRequest.GetProjectId(),
		UserId:    userId,
	})
	require.NoError(t, err)

	t.Run("message from outsider", func(t *testing.T) {
		msg := &entity.Message{
			ProjectID: createChatRequest.GetProjectId(),
			UserID:    userId,
			Content:   gofakeit.Word(),
		}
		err = s.redisClient.SendMessage(ctx, msg)
		require.NoError(t, err)
		msgFake := &entity.Message{
			ProjectID: createChatRequest.GetProjectId(),
			UserID:    userIdFake,
			Content:   gofakeit.Word(),
		}

		err = s.redisClient.SendMessage(ctx, msgFake)
		require.NoError(t, err)
		var resp *api.GetMessagesResponse
		length := 0
		require.Eventuallyf(t, func() bool {
			resp, err = s.grpcClient.GetMessages(ctx, &api.GetMessagesRequest{
				UserId:    userId,
				ProjectId: createChatRequest.GetProjectId(),
				Limit:     10,
				Cursor:    1,
			})
			if resp != nil {
				length = len(resp.Messages)
			}
			return err == nil && len(resp.Messages) == 1
		}, 10*time.Second, 100*time.Millisecond, "expected 1 messages, got %d", length)
		assert.Equal(t, msg.ProjectID, resp.Messages[0].ProjectId)
		assert.Equal(t, msg.UserID, resp.Messages[0].UserId)
		assert.Equal(t, msg.Content, resp.Messages[0].Content)
	})
}
