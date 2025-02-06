package repository

import (
	"chat-service/entity"
	"chat-service/errorz"
	"context"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const prefix = "chat-repository"

func NewChatRepository(client *mongo.Client) *chatRepository {
	return &chatRepository{
		db: client.Database("chat"),
	}
}

type Message struct {
	ID        primitive.ObjectID `bson:"_id,omitempty"`
	ChatID    string             `bson:"chat_id"`
	UserID    string             `bson:"user_id"`
	Message   string             `bson:"message"`
	Timestamp time.Time          `bson:"timestamp"`
}

type chatRepository struct {
	db *mongo.Database
}

func (c *chatRepository) WriteMessage(ctx context.Context, msg *entity.Message) (*entity.Message, error) {
	const op = prefix + ".WriteMessage"
	message := &Message{
		ChatID:    msg.ProjectID,
		UserID:    msg.UserID,
		Message:   msg.Content,
		Timestamp: time.Now(),
	}

	_, err := c.db.Collection("messages").InsertOne(ctx, message)
	if err != nil {
		return nil, errorz.WrapInternal(err, op)
	}
	msg.Time = message.Timestamp.String()
	return msg, nil
}

func (c *chatRepository) GetMessages(ctx context.Context, userID, projectID string, limit, cursor int) ([]*entity.Message, error) {
	const op = prefix + ".GetMessages"
	var messages []*Message
	findOptions := options.Find().
		SetSort(bson.D{{Key: "timestamp", Value: -1}}).
		SetSkip(int64(cursor-1) * int64(limit)).
		SetLimit(int64(limit))

	cur, err := c.db.Collection("messages").Find(ctx, bson.M{"chat_id": projectID}, findOptions)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return []*entity.Message{}, nil
		}
		return nil, errorz.WrapInternal(err, op)
	}

	if err = cur.All(ctx, &messages); err != nil {
		return nil, errorz.WrapInternal(err, op)
	}

	defer cur.Close(ctx)
	var result []*entity.Message
	for _, message := range messages {
		result = append(result, &entity.Message{
			ID:        message.ID.Hex(),
			UserID:    message.UserID,
			ProjectID: message.ChatID,
			Content:   message.Message,
			Time:      message.Timestamp.String(),
		})
	}
	if err := cur.Err(); err != nil {
		return nil, errorz.WrapInternal(err, op)
	}

	return result, nil
}

func (c *chatRepository) CreateChat(ctx context.Context, chat *entity.Chat) (string, error) {
	const op = prefix + ".CreateChat"
	chat.CreatedAt = time.Now()
	result, err := c.db.Collection("chats").InsertOne(ctx, chat)
	if err != nil {
		return "", errorz.WrapInternal(err, op)
	}
	id := result.InsertedID.(primitive.ObjectID)
	return id.Hex(), nil
}

func (c *chatRepository) IsUserInChat(ctx context.Context, projecID string, userID string) (bool, error) {
	const op = prefix + ".IsUserInChat"
	var result bson.M
	err := c.db.Collection("chats").FindOne(ctx, bson.M{"project_id": projecID, "members": userID}).
		Decode(&result)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return false, nil
		}
		return false, errorz.WrapInternal(err, op)
	}
	return true, nil
}

func (c *chatRepository) AddUserToChat(ctx context.Context, userID string, projectID string) error {
	const op = prefix + ".AddUserToChat"
	res, err := c.db.Collection("chats").UpdateOne(ctx, bson.M{"project_id": projectID}, bson.M{"$addToSet": bson.M{"members": userID}})
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return errorz.Wrap(errorz.ErrNotFound, op)
		}
		return errorz.WrapInternal(err, op)
	}
	if res.MatchedCount == 0 {
		return errorz.Wrap(errorz.ErrNotFound, op)
	}
	return nil
}

func (c *chatRepository) RemoveUserFromChat(ctx context.Context, userID string, projectID string) error {
	const op = prefix + ".RemoveUserFromChat"
	_, err := c.db.Collection("chats").UpdateOne(ctx, bson.M{"project_id": projectID}, bson.M{"$pull": bson.M{"members": userID}})
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return errorz.Wrap(errorz.ErrNotFound, op)
		}
		return errorz.WrapInternal(err, op)
	}
	return nil
}

func (c *chatRepository) GetUserChats(ctx context.Context, userID string) ([]string, error) {
	const op = prefix + ".GetUserChats"
	var result []string
	println(userID)
	cur, err := c.db.Collection("chats").Find(ctx, bson.M{"members": userID})
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return result, nil
		}
		return nil, errorz.WrapInternal(err, op)
	}
	defer cur.Close(ctx)

	for cur.Next(ctx) {
		var chat struct {
			ProjectID string `bson:"project_id"`
		}
		if err := cur.Decode(&chat); err != nil {
			return nil, errorz.WrapInternal(err, op)
		}
		result = append(result, chat.ProjectID)
	}

	if err := cur.Err(); err != nil {
		return nil, errorz.WrapInternal(err, op)
	}

	return result, nil
}

func (c *chatRepository) GetChatUsers(ctx context.Context, projectID string) ([]string, error) {
	const op = prefix + ".GetChatUsers"
	var chat entity.Chat
	err := c.db.Collection("chats").FindOne(ctx, bson.M{"project_id": projectID}).Decode(&chat)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return []string{}, nil
		}
		return nil, errorz.WrapInternal(err, op)
	}
	return chat.Members, err
}

func (c *chatRepository) GetChat(ctx context.Context, projectID string) (*entity.Chat, error) {
	const op = prefix + ".GetChat"
	var chat entity.Chat
	err := c.db.Collection("chats").FindOne(ctx, bson.M{"project_id": projectID}).Decode(&chat)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, errorz.Wrap(errorz.ErrNotFound, op)
		}
		return nil, errorz.WrapInternal(err, op)
	}
	return &chat, err
}
