package redisconsumer

import (
	"chat-service/entity"
	"context"
	"encoding/json"
	"log"

	"github.com/redis/go-redis/v9"
)

func (s *Server) processMessage(msg redis.XMessage) {
	message := &entity.Message{
		UserID:    msg.Values["user_id"].(string),
		ProjectID: msg.Values["room_id"].(string),
		Content:   msg.Values["content"].(string),
	}
	s.chatService.WriteMessage(context.TODO(), message)
}

func (s *Server) processSendMessage(ctx context.Context, msg *entity.Message) {
	data, err := json.Marshal(msg)
	if err != nil {
		log.Printf("error: %v", err)
	}

	s.redis.Publish(ctx, s.cfg.PublishChan, data).Err()
}

func (s *Server) StartSendMessagesProcessor(ch chan *entity.Message) {
	for {
		select {
		case msg, ok := <-ch:
			if !ok {
				return
			}
			s.processSendMessage(s.ctx, msg)
		case <-s.ctx.Done():
			return
		}
	}
}
