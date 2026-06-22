package unit_test

import (
	"chat-service/errorz"
	"chat-service/internal/transport/grpc/handlers"
	mocks "chat-service/mocks"
	api "chat-service/pkg/api/chat_v1"
	"context"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func TestHandler_removeUserFromChat(t *testing.T) {
	type mockBehaviour func(s *mocks.MockChatService, userID string, projectID string)

	tests := []struct {
		name           string
		mockBehaviour  mockBehaviour
		ChatID         string
		userID         string
		exceptedStatus int
	}{
		{
			name: "OK",
			mockBehaviour: func(s *mocks.MockChatService, userID string, projectID string) {
				s.EXPECT().RemoveUserFromChat(gomock.Any(), userID, projectID).Return(nil)
			},
			ChatID:         "project-1",
			userID:         "user-1",
			exceptedStatus: int(codes.OK),
		},
		{
			name:           "Empty projectID",
			mockBehaviour:  func(s *mocks.MockChatService, userID string, projectID string) {},
			ChatID:         "",
			userID:         "user-1",
			exceptedStatus: int(codes.InvalidArgument),
		}, {
			name:           "Empty userID",
			mockBehaviour:  func(s *mocks.MockChatService, userID string, projectID string) {},
			ChatID:         "project-1",
			userID:         "",
			exceptedStatus: int(codes.InvalidArgument),
		}, {
			name: "NotFound",
			mockBehaviour: func(s *mocks.MockChatService, userID string, projectID string) {
				s.EXPECT().RemoveUserFromChat(gomock.Any(), userID, projectID).Return(errorz.ErrNotFound)
			},
			ChatID:         "project-1",
			userID:         "user-1",
			exceptedStatus: int(codes.NotFound),
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			repo := mocks.NewMockChatService(ctrl)
			test.mockBehaviour(repo, test.userID, test.ChatID)

			handler := handlers.NewChatHandler(repo)

			resp, err := handler.RemoveUserFromChat(context.Background(),
				&api.RemoveUserFromChatRequest{
					ProjectId: test.ChatID,
					UserId:    test.userID,
				})
			if err == nil {
				require.Equal(t, test.ChatID, resp.GetProjectId())
			}
			if s, ok := status.FromError(err); ok {
				require.Equal(t, test.exceptedStatus, int(s.Code()))
			}
		})
	}
}
