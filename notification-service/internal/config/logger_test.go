package config_test

import (
	"errors"
	"github.com/rs/zerolog"
	"notification-service/internal/config"
	"os"
	"testing"
)

func TestNewLoggerConfig_Success(t *testing.T) {
	// Устанавливаем переменные окружения
	os.Setenv("IS_PRETTY", "true")
	os.Setenv("VERSION", "1.0.0")
	os.Setenv("LOG_LEVEL", "0")
	defer func() {
		os.Unsetenv("IS_PRETTY")
		os.Unsetenv("VERSION")
		os.Unsetenv("LOG_LEVEL")
	}()

	// Создаем конфигурацию
	loggerCfg, err := config.NewLoggerConfig()
	if err != nil {
		t.Fatalf("Expected no error, got: %v", err)
	}

	// Проверяем значения
	if !loggerCfg.IsPretty() {
		t.Errorf("Expected IsPretty to be true, got false")
	}

	if loggerCfg.Version() != "1.0.0" {
		t.Errorf("Expected Version to be '1.0.0', got %v", loggerCfg.Version())
	}

	if loggerCfg.LogLevel() != zerolog.DebugLevel {
		t.Errorf("Expected LogLevel to be 'debug', got %v", loggerCfg.LogLevel())
	}
}

func TestNewLoggerConfig_MissingIsPretty(t *testing.T) {
	// Устанавливаем переменные окружения
	os.Unsetenv("IS_PRETTY")
	os.Setenv("VERSION", "1.0.0")
	os.Setenv("LOG_LEVEL", "0")
	defer func() {
		os.Unsetenv("VERSION")
		os.Unsetenv("LOG_LEVEL")
	}()

	// Создаем конфигурацию
	_, err := config.NewLoggerConfig()
	if !errors.Is(err, config.ErrIsPrettyNotFound) {
		t.Errorf("Expected error '%v', got '%v'", config.ErrIsPrettyNotFound, err)
	}
}

func TestNewLoggerConfig_MissingVersion(t *testing.T) {
	// Устанавливаем переменные окружения
	os.Setenv("IS_PRETTY", "true")
	os.Unsetenv("VERSION")
	os.Setenv("LOG_LEVEL", "0")
	defer func() {
		os.Unsetenv("IS_PRETTY")
		os.Unsetenv("LOG_LEVEL")
	}()

	// Создаем конфигурацию
	_, err := config.NewLoggerConfig()
	if !errors.Is(err, config.ErrVersionNotFound) {
		t.Errorf("Expected error '%v', got '%v'", config.ErrVersionNotFound, err)
	}
}

func TestNewLoggerConfig_InvalidLogLevel(t *testing.T) {
	// Устанавливаем переменные окружения
	os.Setenv("IS_PRETTY", "true")
	os.Setenv("VERSION", "1.0.0")
	os.Setenv("LOG_LEVEL", "invalid")
	defer func() {
		os.Unsetenv("IS_PRETTY")
		os.Unsetenv("VERSION")
		os.Unsetenv("LOG_LEVEL")
	}()

	// Создаем конфигурацию
	_, err := config.NewLoggerConfig()
	if !errors.Is(err, config.ErrLogLevelNotFound) {
		t.Errorf("Expected error '%v', got '%v'", config.ErrLogLevelNotFound, err)
	}
}
