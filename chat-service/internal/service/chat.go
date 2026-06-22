package service

import (
	"chat-service/entity"
	"chat-service/errorz"
	"chat-service/pkg/logger"
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
	"go.uber.org/zap"
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

	DeleteChat(ctx context.Context, projectID string) error
	DeleteMessages(ctx context.Context, projectID string) error
	UpdateChat(ctx context.Context, projectID string, chat *entity.Chat) error

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
	rdb              *redis.Client
}

const prefix = "chat-service"

func NewChatService(chatRepository ChatRepository, UpdatesTransport UpdatesTransport, rdb *redis.Client) *chatService {
	return &chatService{
		repo:             chatRepository,
		updatesTransport: UpdatesTransport,
		sendCh:           make(chan *entity.Message, 10),
		rdb:              rdb,
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
		return nil, errorz.Wrap(errorz.ErrForbidden, "chatService.GetMessages")
	}

	if cursor*limit < 10 {
		cachedMessages, err := s.GetMessagesFromCache(ctx, projectID, limit, cursor)
		if err != nil {
			logger.GetLogger(ctx).Error(ctx, "cache error", zap.Error(err))
		}

		if len(cachedMessages) > 0 {
			if len(cachedMessages) < limit {
				messagesFromDB, dbErr := s.repo.GetMessages(ctx, userID, projectID, limit-len(cachedMessages), cursor+len(cachedMessages))
				if dbErr != nil {
					return cachedMessages, dbErr
				}
				return append(cachedMessages, messagesFromDB...), nil
			}
			return cachedMessages, nil
		}

		messages, err := s.repo.GetMessages(ctx, userID, projectID, 10, 1)
		if err != nil {
			return nil, err
		}
		go func() {
			if err := s.SetMessagesToCache(ctx, projectID, messages, time.Minute); err != nil {
				logger.GetLogger(ctx).Error(ctx, "failed to set cache", zap.Error(err))
			}
		}()
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
		logger.GetLogger(ctx).Error(ctx, "failed to get user in chat", zap.Error(err))
	}
	if !ok {
		return
	}
	msg, err = s.repo.WriteMessage(ctx, msg)
	if err != nil {
		logger.GetLogger(ctx).Error(ctx, "failed to write message", zap.Error(err))
	}
	err = s.CacheMessage(ctx, msg.ProjectID, msg, time.Minute)
	if err != nil {
		logger.GetLogger(ctx).Error(ctx, "failed to set cache", zap.Error(err))
	}
	s.sendCh <- msg
}

func (s *chatService) GetMessagesFromCache(ctx context.Context, projectID string, limit, cursor int) ([]*entity.Message, error) {
	const op = prefix + ".GetMessagesFromCache"
	var messages []*entity.Message

	key := fmt.Sprintf("chat:%s:messages", projectID)
	messagesJSON, err := s.rdb.LRange(ctx, key, int64(cursor*limit), int64(cursor+limit-1)).Result()
	if err != nil {
		if err == redis.Nil {
			return messages, nil
		}
		return messages, errorz.WrapInternal(err, op)
	}
	for _, messageJSON := range messagesJSON {
		var message entity.Message
		if err := json.Unmarshal([]byte(messageJSON), &message); err != nil {
			return messages, errorz.WrapInternal(err, op)
		}
		messages = append(messages, &message)
	}
	return messages, nil
}

func (s *chatService) CacheMessage(ctx context.Context, projectID string, message *entity.Message, ttl time.Duration) error {
	const op = prefix + ".CacheMessages"
	key := fmt.Sprintf("chat:%s:messages", projectID)
	messageJSON, err := json.Marshal(message)
	if err != nil {
		return errorz.WrapInternal(err, op)
	}
	pipe := s.rdb.TxPipeline()
	pipe.LPush(ctx, key, messageJSON)
	pipe.LTrim(ctx, key, 0, 9)
	pipe.Expire(ctx, key, ttl)
	_, err = pipe.Exec(ctx)
	if err != nil {
		return errorz.WrapInternal(err, op)
	}
	return nil
}

func (s *chatService) SetMessagesToCache(ctx context.Context, projectID string, messages []*entity.Message, ttl time.Duration) error {
	const op = prefix + ".SetCacheMessages"
	key := fmt.Sprintf("chat:%s:messages", projectID)
	pipe := s.rdb.TxPipeline()
	for _, message := range messages {
		messageJSON, err := json.Marshal(message)
		if err != nil {
			return errorz.WrapInternal(err, op)
		}
		pipe.LPush(ctx, key, messageJSON)
	}
	pipe.LTrim(ctx, key, 0, 9)
	pipe.Expire(ctx, key, ttl)
	_, err := pipe.Exec(ctx)
	if err != nil {
		return errorz.WrapInternal(err, op)
	}

	return nil
}

func (s *chatService) DeleteChat(ctx context.Context, projectID string) error {
	if err := s.repo.DeleteMessages(ctx, projectID); err != nil {
		return err
	}
	return s.repo.DeleteChat(ctx, projectID)
}

func (s *chatService) UpdateChat(ctx context.Context, projectID string, chat *entity.Chat) error {
	return s.repo.UpdateChat(ctx, projectID, chat)
}
