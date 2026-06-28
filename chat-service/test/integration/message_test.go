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

const (
	readMessageTimeout = 2 * time.Second
	sendMessageTimeout = 200 * time.Millisecond
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

	messages := []*entity.Message{}
	for i := 0; i < 10; i++ {
		msg := &entity.Message{
			ProjectID: createChatRequest.GetProjectId(),
			UserID:    userId,
			Content:   gofakeit.Word(),
		}
		messages = append(messages, msg)
		err := s.redisClient.SendMessage(ctx, msg)
		time.Sleep(sendMessageTimeout)
		require.NoError(t, err)
	}
	time.Sleep(readMessageTimeout)
	resp, err := s.grpcClient.GetMessages(ctx, &api.GetMessagesRequest{
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
		time.Sleep(readMessageTimeout)
		msgFake := &entity.Message{
			ProjectID: createChatRequest.GetProjectId(),
			UserID:    userIdFake,
			Content:   gofakeit.Word(),
		}

		err = s.redisClient.SendMessage(ctx, msgFake)
		require.NoError(t, err)
		resp, err := s.grpcClient.GetMessages(ctx, &api.GetMessagesRequest{
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
	})
}
