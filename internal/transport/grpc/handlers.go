package grpc

import (
	"chat-service/entity"
	"chat-service/errorz"
	api "chat-service/pkg/api/chat_v1"
	"chat-service/pkg/logger"
	"context"

	"go.uber.org/zap"
	"google.golang.org/grpc/codes"
)

type ChatService interface {
	AddUserToChat(ctx context.Context, userID string, projectID string) error
	CreateChat(ctx context.Context, chat *entity.Chat) (string, error)
	GetMessages(ctx context.Context, userID, projectID string, limit, cursor int) ([]*entity.Message, error)
	GetChatUsers(ctx context.Context, projectID string) ([]string, error)
	GetUserChats(ctx context.Context, userID string) ([]string, error)
	GetChat(ctx context.Context, projectID string) (*entity.Chat, error)
	RemoveUserFromChat(ctx context.Context, userID string, projectID string) error
}

func (s *Server) AddUserToChat(ctx context.Context, req *api.AddUserToChatRequest) (*api.AddUserToChatResponse, error) {
	err := s.chatService.AddUserToChat(ctx, req.GetUserId(), req.GetProjectId())
	if err != nil {
		if err := errorz.Parse(err); err != nil {
			if err.Is(errorz.ErrNotFound) {
				err.SetStatusCode(int(codes.NotFound))
			}
		}
		return nil, err
	}
	logger.GetLogger(ctx).Info(ctx, "User added to chat", zap.String("user_id", req.UserId), zap.String("project_id", req.ProjectId))
	return &api.AddUserToChatResponse{ProjectId: req.ProjectId}, nil
}

func (s *Server) CreateChat(ctx context.Context, req *api.CreateChatRequest) (*api.CreateChatResponse, error) {
	chatId, err := s.chatService.CreateChat(ctx, &entity.Chat{
		ProjectID: req.GetProjectId(),
		Members:   req.GetMember(),
		Name:      req.GetName(),
	})
	if err != nil {
		return nil, err
	}
	logger.GetLogger(ctx).Info(ctx, "Chat created", zap.String("chat_id", chatId), zap.String("project_id", req.ProjectId))
	return &api.CreateChatResponse{ChatId: chatId}, nil
}

func (s *Server) GetMessages(ctx context.Context, req *api.GetMessagesRequest) (*api.GetMessagesResponse, error) {
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

func (s *Server) GetChatUsers(ctx context.Context, req *api.GetChatUsersRequest) (*api.GetChatUsersResponse, error) {
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

func (s *Server) GetUserChats(ctx context.Context, req *api.GetUserChatsRequest) (*api.GetUserChatsResponse, error) {
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

func (s *Server) RemoveUserFromChat(ctx context.Context, req *api.RemoveUserFromChatRequest) (*api.RemoveUserFromChatResponse, error) {
	err := s.chatService.RemoveUserFromChat(ctx, req.GetUserId(), req.GetProjectId())
	if err != nil {
		if err := errorz.Parse(err); err != nil {
			if err.Is(errorz.ErrNotFound) {
				err.SetStatusCode(int(codes.NotFound))
			}
		}
		return nil, err
	}
	logger.GetLogger(ctx).Info(ctx, "User removed from chat", zap.String("user_id", req.UserId), zap.String("project_id", req.ProjectId))
	return &api.RemoveUserFromChatResponse{ProjectId: req.ProjectId}, nil
}

func (s *Server) GetChat(ctx context.Context, req *api.GetChatRequest) (*api.GetChatResponse, error) {
	chat, err := s.chatService.GetChat(ctx, req.GetProjectId())
	if err != nil {
		if err := errorz.Parse(err); err != nil {
			if err.Is(errorz.ErrNotFound) {
				err.SetStatusCode(int(codes.NotFound))
			}
		}
		return nil, err
	}
	return &api.GetChatResponse{Chat: &api.Chat{Members: chat.Members, Name: chat.Name, ProjectId: chat.ProjectID, ChatId: chat.GetID()}}, nil
}

// mustEmbedUnimplementedChatServiceServer implements chat_v1.ChatServiceServer.
// Subtle: this method shadows the method (UnsafeChatServiceServer).mustEmbedUnimplementedChatServiceServer of Server.UnsafeChatServiceServer.
func (s *Server) mustEmbedUnimplementedChatServiceServer() {
	panic("unimplemented")
}
