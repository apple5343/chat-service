package handlers

import (
	"chat-service/entity"
	"chat-service/errorz"
	api "chat-service/pkg/api/chat_v1"
	"chat-service/pkg/logger"
	"context"

	"go.uber.org/zap"
	"google.golang.org/grpc/codes"
)

const (
	prefix = "chat-service-grpc"
)

type ChatHandler struct {
	chatService ChatService
	api.UnsafeChatServiceServer
}

// mustEmbedUnimplementedChatServiceServer implements chat_v1.ChatServiceServer.
func (h *ChatHandler) mustEmbedUnimplementedChatServiceServer() {
	panic("unimplemented")
}

func NewChatHandler(chatService ChatService) *ChatHandler {
	return &ChatHandler{chatService: chatService}
}

type ChatService interface {
	AddUserToChat(ctx context.Context, userID string, projectID string) error
	CreateChat(ctx context.Context, chat *entity.Chat) (string, error)
	GetMessages(ctx context.Context, userID, projectID string, limit, cursor int) ([]*entity.Message, error)
	GetChatUsers(ctx context.Context, projectID string) ([]string, error)
	GetUserChats(ctx context.Context, userID string) ([]string, error)
	GetChat(ctx context.Context, projectID string) (*entity.Chat, error)
	RemoveUserFromChat(ctx context.Context, userID string, projectID string) error
	DeleteChat(ctx context.Context, projectID string) error
	UpdateChat(ctx context.Context, projectID string, chat *entity.Chat) error
}

func (h *ChatHandler) AddUserToChat(ctx context.Context, req *api.AddUserToChatRequest) (*api.AddUserToChatResponse, error) {
	if req.GetUserId() == "" || req.GetProjectId() == "" {
		return nil, errorz.BadRequest()
	}
	err := h.chatService.AddUserToChat(ctx, req.GetUserId(), req.GetProjectId())
	if err != nil {
		if err := errorz.Parse(err); err != nil {
			if err.Is(errorz.ErrNotFound) {
				err.SetStatusCode(int(codes.NotFound))
			}
		}
		return nil, err
	}
	if logger := logger.GetLogger(ctx); logger != nil {
		logger.Info(ctx, "User added to chat", zap.String("user_id", req.UserId), zap.String("project_id", req.ProjectId))
	}
	return &api.AddUserToChatResponse{ProjectId: req.ProjectId}, nil
}

func (s *ChatHandler) CreateChat(ctx context.Context, req *api.CreateChatRequest) (*api.CreateChatResponse, error) {
	if req.GetProjectId() == "" || req.GetName() == "" {
		return nil, errorz.BadRequest()
	}
	chatId, err := s.chatService.CreateChat(ctx, &entity.Chat{
		ProjectID: req.GetProjectId(),
		Members:   req.GetMember(),
		Name:      req.GetName(),
	})
	if err != nil {
		return nil, err
	}
	if logger := logger.GetLogger(ctx); logger != nil {
		logger.Info(ctx, "Chat created", zap.String("chat_id", chatId), zap.String("project_id", req.ProjectId))
	}
	return &api.CreateChatResponse{ChatId: chatId}, nil
}

func (s *ChatHandler) GetMessages(ctx context.Context, req *api.GetMessagesRequest) (*api.GetMessagesResponse, error) {
	if req.GetUserId() == "" || req.GetProjectId() == "" || req.GetLimit() <= 0 || req.GetCursor() <= 0 {
		return nil, errorz.BadRequest()
	}
	messages, err := s.chatService.GetMessages(ctx, req.GetUserId(), req.GetProjectId(), int(req.Limit), int(req.Cursor))
	if err != nil {
		if err := errorz.Parse(err); err != nil {
			if err.Is(errorz.ErrNotFound) {
				err.SetStatusCode(int(codes.NotFound))
			} else if err.Is(errorz.ErrForbidden) {
				err.SetStatusCode(int(codes.PermissionDenied))
			}
		}
		return nil, err
	}
	result := make([]*api.Message, 0, len(messages))
	for _, message := range messages {
		result = append(result, &api.Message{
			MessageId: message.ID,
			ProjectId: message.ProjectID,
			UserId:    message.UserID,
			Content:   message.Content,
			Timestamp: message.Time,
		})
	}
	return &api.GetMessagesResponse{Messages: result}, nil
}

func (s *ChatHandler) GetChatUsers(ctx context.Context, req *api.GetChatUsersRequest) (*api.GetChatUsersResponse, error) {
	if req.GetProjectId() == "" {
		return nil, errorz.BadRequest()
	}
	users, err := s.chatService.GetChatUsers(ctx, req.GetProjectId())
	if err != nil {
		if err := errorz.Parse(err); err != nil {
			if err.Is(errorz.ErrNotFound) {
				err.SetStatusCode(int(codes.NotFound))
			}
		}
		return nil, err
	}
	return &api.GetChatUsersResponse{UserIds: users}, nil

}

func (s *ChatHandler) GetUserChats(ctx context.Context, req *api.GetUserChatsRequest) (*api.GetUserChatsResponse, error) {
	if req.GetUserId() == "" {
		return nil, errorz.BadRequest()
	}
	chats, err := s.chatService.GetUserChats(ctx, req.GetUserId())
	if err != nil {
		if err := errorz.Parse(err); err != nil {
			if err.Is(errorz.ErrNotFound) {
				err.SetStatusCode(int(codes.NotFound))
			}
		}
		return nil, err
	}
	return &api.GetUserChatsResponse{ProjectIds: chats}, nil
}

func (s *ChatHandler) RemoveUserFromChat(ctx context.Context, req *api.RemoveUserFromChatRequest) (*api.RemoveUserFromChatResponse, error) {
	if req.GetUserId() == "" || req.GetProjectId() == "" {
		return nil, errorz.BadRequest()
	}
	err := s.chatService.RemoveUserFromChat(ctx, req.GetUserId(), req.GetProjectId())
	if err != nil {
		if err := errorz.Parse(err); err != nil {
			if err.Is(errorz.ErrNotFound) {
				err.SetStatusCode(int(codes.NotFound))
			}
		}
		return nil, err
	}
	if logger := logger.GetLogger(ctx); logger != nil {
		logger.Info(ctx, "User removed from chat", zap.String("user_id", req.UserId), zap.String("project_id", req.ProjectId))
	}
	return &api.RemoveUserFromChatResponse{ProjectId: req.ProjectId}, nil
}

func (s *ChatHandler) GetChat(ctx context.Context, req *api.GetChatRequest) (*api.GetChatResponse, error) {
	if req.GetProjectId() == "" {
		return nil, errorz.BadRequest()
	}
	chat, err := s.chatService.GetChat(ctx, req.GetProjectId())
	if err != nil {
		if err := errorz.Parse(err); err != nil {
			if err.Is(errorz.ErrNotFound) {
				err.SetStatusCode(int(codes.NotFound))
			}
		}
		return nil, err
	}
	return &api.GetChatResponse{Chat: &api.Chat{Members: chat.Members, Name: chat.Name,
		ProjectId: chat.ProjectID, ChatId: chat.GetID(), CreatedAt: chat.CreatedAt.String()}}, nil
}

func (s *ChatHandler) DeleteChat(ctx context.Context, req *api.DeleteChatRequest) (*api.DeleteChatResponse, error) {
	if req.GetProjectId() == "" {
		return nil, errorz.BadRequest()
	}
	err := s.chatService.DeleteChat(ctx, req.GetProjectId())
	if err != nil {
		if err := errorz.Parse(err); err != nil {
			if err.Is(errorz.ErrNotFound) {
				err.SetStatusCode(int(codes.NotFound))
			}
		}
		return nil, err
	}
	if logger := logger.GetLogger(ctx); logger != nil {
		logger.Info(ctx, "Chat deleted", zap.String("project_id", req.ProjectId))
	}
	return &api.DeleteChatResponse{ProjectId: req.ProjectId}, nil
}

func (s *ChatHandler) UpdateChat(ctx context.Context, req *api.UpdateChatRequest) (*api.UpdateChatResponse, error) {
	reqChat := req.GetChat()
	if reqChat.GetName() == "" || reqChat.GetProjectId() == "" {
		return nil, errorz.BadRequest()
	}
	chat := &entity.Chat{
		Name:    reqChat.GetName(),
		Members: reqChat.GetMembers(),
	}
	err := s.chatService.UpdateChat(ctx, req.GetChat().GetChatId(), chat)
	if err != nil {
		if err := errorz.Parse(err); err != nil {
			if err.Is(errorz.ErrNotFound) {
				err.SetStatusCode(int(codes.NotFound))
			}
		}
		return nil, err
	}
	if logger := logger.GetLogger(ctx); logger != nil {
		logger.Info(ctx, "Chat updated", zap.String("project_id", reqChat.GetProjectId()))
	}
	return &api.UpdateChatResponse{Chat: &api.Chat{
		Name:    reqChat.GetName(),
		Members: reqChat.GetMembers(),
	}}, nil
}
