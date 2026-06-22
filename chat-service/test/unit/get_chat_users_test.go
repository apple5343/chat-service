package unit_test

import (
	"chat-service/errorz"
	"chat-service/internal/transport/grpc/handlers"
	mocks "chat-service/mocks"
	api "chat-service/pkg/api/chat_v1"
	"context"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func TestHandler_getChatUsers(t *testing.T) {
	type mockBehaviour func(s *mocks.MockChatService, projectID string)

	tests := []struct {
		name           string
		mockBehaviour  mockBehaviour
		ChatID         string
		ChatUsersIds   []string
		exceptedStatus int
	}{
		{
			name: "OK",
			mockBehaviour: func(s *mocks.MockChatService, projectID string) {
				s.EXPECT().GetChatUsers(gomock.Any(), projectID).Return([]string{"user-1", "user-2"}, nil)
			},
			ChatID:         "project-1",
			ChatUsersIds:   []string{"user-1", "user-2"},
			exceptedStatus: int(codes.OK),
		},
		{
			name:           "Empty projectID",
			mockBehaviour:  func(s *mocks.MockChatService, projectID string) {},
			ChatID:         "",
			ChatUsersIds:   nil,
			exceptedStatus: int(codes.InvalidArgument),
		},
		{
			name: "Not found",
			mockBehaviour: func(s *mocks.MockChatService, projectID string) {
				s.EXPECT().GetChatUsers(gomock.Any(), projectID).Return(nil, errorz.ErrNotFound)
			},
			ChatID:         "project-1",
			ChatUsersIds:   nil,
			exceptedStatus: int(codes.NotFound),
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			service := mocks.NewMockChatService(ctrl)
			test.mockBehaviour(service, test.ChatID)

			handler := handlers.NewChatHandler(service)
			users, err := handler.GetChatUsers(context.Background(), &api.GetChatUsersRequest{ProjectId: test.ChatID})
			if err == nil {
				assert.Len(t, users.GetUserIds(), len(test.ChatUsersIds))
				assert.Equal(t, test.ChatUsersIds, users.GetUserIds())
			}

			if s, ok := status.FromError(err); ok {
				assert.Equal(t, test.exceptedStatus, int(s.Code()))
			}
		})
	}
}
