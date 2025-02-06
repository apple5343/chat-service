package entity

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Message struct {
	ID        string `json:"message_id,omitempty" bson:"_id"`
	Content   string `json:"content" bson:"content"`
	ProjectID string `json:"room_id" bson:"room_id"`
	UserID    string `json:"user_id" bson:"user_id"`
	Time      string `json:"time,omitempty" bson:"time"`
}

type Chat struct {
	ID        primitive.ObjectID `bson:"_id,omitempty"`
	ProjectID string             `bson:"project_id"`
	Name      string             `bson:"name"`
	Members   []string           `bson:"members"`
	CreatedAt time.Time          `bson:"created_at"`
}

func (c *Chat) GetID() string {
	return c.ID.Hex()
}
