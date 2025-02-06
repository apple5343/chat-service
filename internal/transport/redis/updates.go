package redisconsumer

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/redis/go-redis/v9"
)

const (
	updatesChannel = "updates-channel"
)

type updatesTransport struct {
	redis *redis.Client
}

func NewUpdatesTransport(redis *redis.Client) *updatesTransport {
	return &updatesTransport{
		redis: redis,
	}
}

func (t *updatesTransport) AddUserToChat(ctx context.Context, userID string, projectID string) error {
	var msg struct {
		Type   string `json:"type"`
		UserID string `json:"user_id"`
		RoomID string `json:"room_id"`
	}

	msg.UserID = userID
	msg.RoomID = projectID
	msg.Type = "user_added"
	data, err := json.Marshal(msg)
	if err != nil {
		return err
	}
	fmt.Println(string(data))
	return t.redis.Publish(ctx, updatesChannel, data).Err()
}

func (u *updatesTransport) RemoveUserFromChat(ctx context.Context, userID string, projectID string) error {
	var msg struct {
		Type   string `json:"type"`
		UserID string `json:"user_id"`
		RoomID string `json:"room_id"`
	}

	msg.UserID = userID
	msg.RoomID = projectID
	msg.Type = "user_removed"
	data, err := json.Marshal(msg)
	if err != nil {
		return err
	}
	return u.redis.Publish(ctx, updatesChannel, data).Err()
}
