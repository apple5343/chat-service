package service

import (
	"api-gateway/internal/entity"
	"context"
	"encoding/json"
	"log"

	"github.com/redis/go-redis/v9"
)

const (
	messagesStream  = "chat-stream"
	messagesChannel = "chat-channel"
	updatesChannel  = "updates-channel"
)

type Chat struct {
	ctx context.Context
	rdb *redis.Client
}

func NewChat(rdb *redis.Client) *Chat {
	return &Chat{
		ctx: context.Background(),
		rdb: rdb,
	}
}

func (c *Chat) WriteMessage(msg *entity.Message) {
	_, err := c.rdb.XAdd(c.ctx, &redis.XAddArgs{
		Stream: messagesStream,
		Values: map[string]string{
			"content": msg.Content,
			"user_id": msg.UserID,
			"room_id": msg.RoomID,
		},
	}).Result()
	if err != nil {
		log.Printf("error: %v", err)
	}
}

func (c *Chat) ReadMessage(ch chan *entity.Message) {
	subscriber := c.rdb.Subscribe(c.ctx, messagesChannel)

	messages := subscriber.Channel()
	go func() {
		for msg := range messages {
			var message entity.Message
			err := json.Unmarshal([]byte(msg.Payload), &message)
			message.Type = "message"
			if err != nil {
				log.Printf("error: %v", err)
				continue
			}
			ch <- &message
		}
	}()
}

func (c *Chat) ReadUpdates(ch chan *entity.Update) {
	subscriber := c.rdb.Subscribe(c.ctx, updatesChannel)

	updates := subscriber.Channel()
	go func() {
		for msg := range updates {
			var update entity.Update
			err := json.Unmarshal([]byte(msg.Payload), &update)
			if err != nil {
				log.Printf("error: %v", err)
				continue
			}
			update.Data = msg.Payload
			ch <- &update
		}
	}()
}
