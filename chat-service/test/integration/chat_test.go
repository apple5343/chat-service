package integration

import (
	"testing"

	api "chat-service/pkg/api/chat_v1"
	"chat-service/test/integration/grpc"
	"chat-service/test/integration/redis"

	"github.com/brianvoe/gofakeit/v6"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type ChatSuite struct {
	suite.Suite
	base        *BaseTestSuite
	grpcClient  api.ChatServiceClient
	redisClient *redis.RedisClient
}

func (s *ChatSuite) SetupSuite() {
	s.grpcClient = grpc.NewClient(s.T(), s.base.config.GrpcConfig)
	s.redisClient = redis.NewClient(s.T(), s.base.config.RedisConfig)
}

func (s *ChatSuite) TearDownSuite() {
	s.redisClient.Close()
}

func (s *ChatSuite) TestCreateChat() {
	s.T().Run("happy path", s.CreateChat)
	s.T().Run("invalid", s.CreateChatInvalid)
}

func (s *ChatSuite) CreateChat(t *testing.T) {
	ctx := s.base.ctx

	req := api.CreateChatRequest{
		ProjectId: gofakeit.UUID(),
		Name:      gofakeit.Name(),
		Member:    []string{gofakeit.UUID()},
	}

	_, err := s.grpcClient.CreateChat(ctx, &req)
	require.NoError(t, err)

	resp, err := s.grpcClient.GetChat(ctx, &api.GetChatRequest{
		ProjectId: req.ProjectId,
	})
	require.NoError(t, err)
	require.NotNil(t, resp)

	assert.Equal(t, req.ProjectId, resp.Chat.ProjectId)
	assert.Equal(t, req.Name, resp.Chat.Name)
	assert.Equal(t, req.Member, resp.Chat.Members)
}

func (s *ChatSuite) CreateChatInvalid(t *testing.T) {
	ctx := s.base.ctx

	t.Run("empty project id", func(t *testing.T) {
		_, err := s.grpcClient.CreateChat(ctx, &api.CreateChatRequest{})
		require.Error(t, err)
		assert.Contains(t, err.Error(), "rpc error: code = InvalidArgument desc = bad request")
	})

	t.Run("empty name", func(t *testing.T) {
		_, err := s.grpcClient.CreateChat(ctx, &api.CreateChatRequest{
			ProjectId: gofakeit.UUID(),
		})
		require.Error(t, err)
		assert.Contains(t, err.Error(), "rpc error: code = InvalidArgument desc = bad request")
	})
}

func (s *ChatSuite) GetChatHistory(t *testing.T) {
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
	resp, err := s.grpcClient.GetMessages(ctx, &api.GetMessagesRequest{
		UserId:    userId,
		ProjectId: createChatRequest.GetProjectId(),
		Limit:     10,
		Cursor:    1,
	})
	require.NoError(t, err)
	require.NotNil(t, resp)
	require.Equal(t, 1, len(resp.Messages))
}

func (s *ChatSuite) GetChatHistoryInvalid(t *testing.T) {
	ctx := s.base.ctx

	createChatRequest := api.CreateChatRequest{
		ProjectId: gofakeit.UUID(),
		Name:      gofakeit.Name(),
		Member:    []string{},
	}

	_, err := s.grpcClient.CreateChat(ctx, &createChatRequest)
	require.NoError(t, err)

	t.Run("permission denied", func(t *testing.T) {
		_, err := s.grpcClient.GetMessages(ctx, &api.GetMessagesRequest{
			UserId:    gofakeit.UUID(),
			ProjectId: createChatRequest.GetProjectId(),
			Limit:     10,
			Cursor:    1,
		})
		require.Error(t, err)
		if s, ok := status.FromError(err); ok {
			assert.EqualValues(t, codes.PermissionDenied, s.Code())
		}
	})
}
