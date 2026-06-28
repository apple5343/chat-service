package integration

import (
	api "chat-service/pkg/api/chat_v1"
	"chat-service/test/integration/grpc"
	"chat-service/test/integration/redis"
	"testing"

	"github.com/brianvoe/gofakeit/v6"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
)

type UserSuite struct {
	suite.Suite
	base        *BaseTestSuite
	grpcClient  api.ChatServiceClient
	redisClient *redis.RedisClient
}

func (s *UserSuite) SetupSuite() {
	s.grpcClient = grpc.NewClient(s.T(), s.base.config.GrpcConfig)
	s.redisClient = redis.NewClient(s.T(), s.base.config.RedisConfig)
}

func (s *UserSuite) TearDownSuite() {
	s.redisClient.Close()
}

func (s *UserSuite) TestAddUser() {
	s.T().Run("happy path", s.AddUser)
	s.T().Run("invalid", s.AddUserInvalid)
}

func (s *UserSuite) AddUser(t *testing.T) {
	ctx := s.base.ctx

	createChatRequest := api.CreateChatRequest{
		ProjectId: gofakeit.UUID(),
		Name:      gofakeit.Name(),
		Member:    []string{},
	}

	_, err := s.grpcClient.CreateChat(ctx, &createChatRequest)
	require.NoError(t, err)

	members := []string{}
	for i := 0; i < 10; i++ {
		id := gofakeit.UUID()
		members = append(members, id)
		r, err := s.grpcClient.AddUserToChat(ctx, &api.AddUserToChatRequest{
			ProjectId: createChatRequest.GetProjectId(),
			UserId:    id,
		})
		require.NoError(t, err)
		require.Equal(t, createChatRequest.GetProjectId(), r.GetProjectId())
	}

	resp, err := s.grpcClient.GetChat(ctx, &api.GetChatRequest{
		ProjectId: createChatRequest.ProjectId,
	})
	require.NoError(t, err)
	require.NotNil(t, resp)

	assert.Equal(t, createChatRequest.ProjectId, resp.Chat.GetProjectId())
	assert.Equal(t, createChatRequest.Name, resp.Chat.GetName())
	assert.Equal(t, members, resp.Chat.GetMembers())
}

func (s *UserSuite) AddUserInvalid(t *testing.T) {
	t.Run("not existing chat", func(t *testing.T) {
		ctx := s.base.ctx
		_, err := s.grpcClient.AddUserToChat(ctx, &api.AddUserToChatRequest{
			ProjectId: gofakeit.UUID(),
			UserId:    gofakeit.UUID(),
		})
		require.Error(t, err)
	})
}
