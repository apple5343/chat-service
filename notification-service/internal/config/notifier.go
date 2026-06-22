package config

import (
	"errors"
	"os"
)

var (
	ErrBotTokenNotFound = errors.New("bot token not found")
)

const (
	botTokenEnvName = "BOT_TOKEN"
)

type NotifierConfig interface {
	BotToken() string
}

type notifierConfig struct {
	botToken string
}

func NewNotifierConfig() (NotifierConfig, error) {
	botToken := os.Getenv(botTokenEnvName)
	if len(botToken) == 0 {
		return nil, ErrBotTokenNotFound
	}

	return &notifierConfig{
		botToken: botToken,
	}, nil
}

func (cfg *notifierConfig) BotToken() string {
	return cfg.botToken
}
