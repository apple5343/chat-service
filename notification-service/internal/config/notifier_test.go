package config_test

import (
	"errors"
	"notification-service/internal/config"
	"os"
	"testing"
)

func TestNewNotifierConfig_Success(t *testing.T) {
	// Устанавливаем переменные окружения
	os.Setenv("BOT_TOKEN", "sample-bot-token")
	defer os.Unsetenv("BOT_TOKEN")

	// Создаем конфигурацию
	notifierCfg, err := config.NewNotifierConfig()
	if err != nil {
		t.Fatalf("Expected no error, got: %v", err)
	}

	// Проверяем значение botToken
	if notifierCfg.BotToken() != "sample-bot-token" {
		t.Errorf("Expected BotToken to be 'sample-bot-token', got '%v'", notifierCfg.BotToken())
	}
}

func TestNewNotifierConfig_MissingBotToken(t *testing.T) {
	// Устанавливаем переменные окружения
	os.Unsetenv("BOT_TOKEN")

	// Создаем конфигурацию
	_, err := config.NewNotifierConfig()
	if !errors.Is(err, config.ErrBotTokenNotFound) {
		t.Errorf("Expected error '%v', got '%v'", config.ErrBotTokenNotFound, err)
	}
}
