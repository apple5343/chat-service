package service

import (
	"chat-service/entity"
	"chat-service/errorz"
	"context"
	"log"
)

type ChatRepository interface {
	GetChatUsers(ctx context.Context, projectID string) ([]string, error)
	GetUserChats(ctx context.Context, userID string) ([]string, error)
	GetMessages(ctx context.Context, userID, projectID string, limit, cursor int) ([]*entity.Message, error)
	GetChat(ctx context.Context, projectID string) (*entity.Chat, error)

	AddUserToChat(ctx context.Context, projecID string, userID string) error
	WriteMessage(ctx context.Context, msg *entity.Message) (*entity.Message, error)
	CreateChat(ctx context.Context, chat *entity.Chat) (string, error)

	RemoveUserFromChat(ctx context.Context, projecID string, userID string) error

	IsUserInChat(ctx context.Context, projecID string, userID string) (bool, error)
}

type UpdatesTransport interface {
	AddUserToChat(ctx context.Context, userID string, projectID string) error
	RemoveUserFromChat(ctx context.Context, userID string, projectID string) error
}

type chatService struct {
	repo             ChatRepository
	updatesTransport UpdatesTransport
	sendCh           chan *entity.Message
}

func NewChatService(chatRepository ChatRepository, UpdatesTransport UpdatesTransport) *chatService {
	return &chatService{
		repo:             chatRepository,
		updatesTransport: UpdatesTransport,
		sendCh:           make(chan *entity.Message, 10),
	}
}

func (s *chatService) AddUserToChat(ctx context.Context, userID string, projectID string) error {
	if err := s.repo.AddUserToChat(ctx, userID, projectID); err != nil {
		return err
	}
	return s.updatesTransport.AddUserToChat(ctx, userID, projectID)
}

func (s *chatService) GetMessages(ctx context.Context, userID, projectID string, limit, cursor int) ([]*entity.Message, error) {
	ok, err := s.repo.IsUserInChat(ctx, projectID, userID)
	if err != nil {
		return nil, err
	}
	if !ok {
		return nil, errorz.ErrForbidden
	}
	return s.repo.GetMessages(ctx, userID, projectID, limit, cursor)

}

func (s *chatService) CreateChat(ctx context.Context, chat *entity.Chat) (string, error) {
	result, err := s.repo.CreateChat(ctx, chat)
	if err != nil {
		return "", err
	}
	return result, nil
}

func (s *chatService) GetChatUsers(ctx context.Context, projectID string) ([]string, error) {
	return s.repo.GetChatUsers(ctx, projectID)
}

func (s *chatService) GetUserChats(ctx context.Context, userID string) ([]string, error) {
	return s.repo.GetUserChats(ctx, userID)
}

func (s *chatService) RemoveUserFromChat(ctx context.Context, userID string, projectID string) error {
	if err := s.repo.RemoveUserFromChat(ctx, userID, projectID); err != nil {
		return err
	}
	return s.updatesTransport.RemoveUserFromChat(ctx, userID, projectID)
}

func (s *chatService) GetChat(ctx context.Context, projectID string) (*entity.Chat, error) {
	return s.repo.GetChat(ctx, projectID)
}

func (s *chatService) ReadMessage(ch chan *entity.Message) {
	go func() {
		for msg := range s.sendCh {
			ch <- msg
		}
	}()
}

func (s *chatService) WriteMessage(ctx context.Context, msg *entity.Message) {
	ok, err := s.repo.IsUserInChat(ctx, msg.ProjectID, msg.UserID)
	if err != nil {
		log.Printf("error: %v", err)
	}
	if !ok {
		return
	}
	msg, err = s.repo.WriteMessage(ctx, msg)
	if err != nil {
		log.Printf("error: %v", err)
	}
	s.sendCh <- msg
}
