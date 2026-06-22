package unit_test

import (
	"chat-service/entity"
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

func TestHandler_getMessages(t *testing.T) {
	type mockBehaviour func(s *mocks.MockChatService, chatID string, userID string, limit int, cursor int)

	msg := &entity.Message{
		UserID:    "user-1",
		ProjectID: "project-1",
		Content:   "hello",
	}
	messages := []*entity.Message{msg, msg, msg, msg, msg}

	tests := []struct {
		name           string
		mockBehaviour  mockBehaviour
		ChatID         string
		userID         string
		limit          int
		cursor         int
		exceptedStatus int
		exceptedMsgs   []*entity.Message
	}{
		{
			name: "OK",
			mockBehaviour: func(s *mocks.MockChatService, chatID string, userID string, limit int, cursor int) {
				s.EXPECT().GetMessages(gomock.Any(), userID, chatID, limit, cursor).Return(messages, nil)
			},
			ChatID:         "project-1",
			userID:         "user-1",
			limit:          5,
			cursor:         1,
			exceptedStatus: int(codes.OK),
			exceptedMsgs:   messages,
		},
		{
			name:           "Empty projectID",
			mockBehaviour:  func(s *mocks.MockChatService, chatID string, userID string, limit int, cursor int) {},
			ChatID:         "",
			userID:         "user-1",
			limit:          5,
			cursor:         1,
			exceptedStatus: int(codes.InvalidArgument),
		}, {
			name:           "Empty userID",
			mockBehaviour:  func(s *mocks.MockChatService, chatID string, userID string, limit int, cursor int) {},
			ChatID:         "project-1",
			userID:         "",
			limit:          5,
			cursor:         1,
			exceptedStatus: int(codes.InvalidArgument),
		}, {
			name:           "Invalid limit",
			mockBehaviour:  func(s *mocks.MockChatService, chatID string, userID string, limit int, cursor int) {},
			ChatID:         "project-1",
			userID:         "user-1",
			limit:          -1,
			cursor:         1,
			exceptedStatus: int(codes.InvalidArgument),
		}, {
			name:           "Invalid cursor",
			mockBehaviour:  func(s *mocks.MockChatService, chatID string, userID string, limit int, cursor int) {},
			ChatID:         "project-1",
			userID:         "user-1",
			limit:          5,
			cursor:         -1,
			exceptedStatus: int(codes.InvalidArgument),
		}, {
			name: "Forbidden",
			mockBehaviour: func(s *mocks.MockChatService, chatID string, userID string, limit int, cursor int) {
				s.EXPECT().GetMessages(gomock.Any(), userID, chatID, limit, cursor).Return(nil, errorz.ErrForbidden)
			},
			ChatID:         "project-1",
			userID:         "user-1",
			limit:          5,
			cursor:         1,
			exceptedStatus: int(codes.PermissionDenied),
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			chatService := mocks.NewMockChatService(ctrl)
			test.mockBehaviour(chatService, test.ChatID, test.userID, test.limit, test.cursor)

			handler := handlers.NewChatHandler(chatService)

			resp, err := handler.GetMessages(context.Background(), &api.GetMessagesRequest{
				ProjectId: test.ChatID,
				UserId:    test.userID,
				Limit:     int32(test.limit),
				Cursor:    int32(test.cursor),
			})

			if err == nil {
				assert.Len(t, resp.Messages, len(test.exceptedMsgs))
				for i := 0; i < len(test.exceptedMsgs); i++ {
					assert.Equal(t, test.exceptedMsgs[i].Content, resp.Messages[i].Content)
					assert.Equal(t, test.exceptedMsgs[i].UserID, resp.Messages[i].UserId)
					assert.Equal(t, test.exceptedMsgs[i].ProjectID, resp.Messages[i].ProjectId)
				}
			}
			if s, ok := status.FromError(err); ok {
				assert.Equal(t, test.exceptedStatus, int(s.Code()))
			}
		})
	}
}
