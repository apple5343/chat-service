package unit_test

import (
	"chat-service/errorz"
	"chat-service/internal/transport/grpc/handlers"
	mocks "chat-service/mocks"
	api "chat-service/pkg/api/chat_v1"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func TestHandler_getUserChats(t *testing.T) {
	type mockBehaviour func(s *mocks.MockChatService, userID string, chatIds []string)

	tests := []struct {
		name           string
		mockBehaviour  mockBehaviour
		UserID         string
		ChatIds        []string
		exceptedStatus int
	}{
		{
			name: "OK",
			mockBehaviour: func(s *mocks.MockChatService, userID string, chatIds []string) {
				s.EXPECT().GetUserChats(gomock.Any(), userID).Return(chatIds, nil)
			},
			UserID:         "user-1",
			ChatIds:        []string{"project-1", "project-2"},
			exceptedStatus: int(codes.OK),
		},
		{
			name:           "Empty userID",
			mockBehaviour:  func(s *mocks.MockChatService, userID string, chatIds []string) {},
			UserID:         "",
			ChatIds:        nil,
			exceptedStatus: int(codes.InvalidArgument),
		},
		{
			name: "Not found",
			mockBehaviour: func(s *mocks.MockChatService, userID string, chatIds []string) {
				s.EXPECT().GetUserChats(gomock.Any(), userID).Return(nil, errorz.ErrNotFound)
			},
			UserID:         "user-1",
			ChatIds:        nil,
			exceptedStatus: int(codes.NotFound),
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			chatService := mocks.NewMockChatService(ctrl)
			test.mockBehaviour(chatService, test.UserID, test.ChatIds)

			handler := handlers.NewChatHandler(chatService)
			resp, err := handler.GetUserChats(nil, &api.GetUserChatsRequest{UserId: test.UserID})
			if err == nil {
				assert.Equal(t, test.ChatIds, resp.GetProjectIds())
			}

			if s, ok := status.FromError(err); ok {
				assert.Equal(t, test.exceptedStatus, int(s.Code()))
			}
		})
	}
}
