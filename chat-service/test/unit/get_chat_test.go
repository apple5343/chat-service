package unit_test

import (
	"chat-service/entity"
	"chat-service/errorz"
	"chat-service/internal/transport/grpc/handlers"
	mocks "chat-service/mocks"
	api "chat-service/pkg/api/chat_v1"
	"testing"
	"time"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func TestHandler_getChat(t *testing.T) {
	type mockBehaviour func(s *mocks.MockChatService, projectID string)
	chat := &entity.Chat{ProjectID: "project-1", Name: "project-1", CreatedAt: time.Now()}

	tests := []struct {
		name           string
		mockBehaviour  mockBehaviour
		ChatID         string
		Chat           *entity.Chat
		exceptedStatus int
	}{
		{
			name: "OK",
			mockBehaviour: func(s *mocks.MockChatService, projectID string) {
				s.EXPECT().GetChat(gomock.Any(), projectID).Return(chat, nil)
			},
			ChatID:         "project-1",
			Chat:           chat,
			exceptedStatus: int(codes.OK),
		},
		{
			name:           "Empty projectID",
			mockBehaviour:  func(s *mocks.MockChatService, projectID string) {},
			ChatID:         "",
			Chat:           nil,
			exceptedStatus: int(codes.InvalidArgument),
		},
		{
			name: "Not found",
			mockBehaviour: func(s *mocks.MockChatService, projectID string) {
				s.EXPECT().GetChat(gomock.Any(), projectID).Return(nil, errorz.ErrNotFound)
			},
			ChatID:         "project-1",
			Chat:           nil,
			exceptedStatus: int(codes.NotFound),
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mockChatService := mocks.NewMockChatService(ctrl)
			test.mockBehaviour(mockChatService, test.ChatID)

			handler := handlers.NewChatHandler(mockChatService)

			resp, err := handler.GetChat(nil, &api.GetChatRequest{ProjectId: test.ChatID})
			if err == nil {
				assert.Equal(t, test.Chat.ProjectID, resp.GetChat().GetProjectId())
				assert.Equal(t, test.Chat.Name, resp.GetChat().GetName())
				assert.Equal(t, test.Chat.CreatedAt.String(), resp.GetChat().GetCreatedAt())
			}

			if s, ok := status.FromError(err); ok {
				assert.Equal(t, test.exceptedStatus, int(s.Code()))
			}
		})
	}
}
