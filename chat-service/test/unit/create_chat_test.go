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
)

func TestHandler_createChat(t *testing.T) {
	type mockBehaviour func(s *mocks.MockChatService, chat *entity.Chat)

	tests := []struct {
		name           string
		mockBehaviour  mockBehaviour
		chat           *entity.Chat
		exceptedStatus int
	}{
		{
			name: "OK",
			mockBehaviour: func(s *mocks.MockChatService, chat *entity.Chat) {
				s.EXPECT().CreateChat(gomock.Any(), chat).Return(chat.GetID(), nil)
			},
			chat: &entity.Chat{
				ProjectID: "project-1",
				Name:      "some_name",
				Members:   []string{"user-1", "user-2"},
			},
			exceptedStatus: int(codes.OK),
		},
		{
			name:          "Invalid request",
			mockBehaviour: func(s *mocks.MockChatService, chat *entity.Chat) {},
			chat: &entity.Chat{
				ProjectID: "",
				Name:      "some_name",
				Members:   []string{"user-1", "user-2"},
			},
			exceptedStatus: int(codes.InvalidArgument),
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			repo := mocks.NewMockChatService(ctrl)
			test.mockBehaviour(repo, test.chat)

			handler := handlers.NewChatHandler(repo)

			req := &api.CreateChatRequest{
				ProjectId: test.chat.ProjectID,
				Name:      test.chat.Name,
				Member:    test.chat.Members,
			}

			_, err := handler.CreateChat(context.Background(), req)

			if customError := errorz.Parse(err); customError != nil {
				assert.Equal(t, test.exceptedStatus, customError.StatusCode())
			}
		})
	}
}
