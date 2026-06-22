package config

import (
	"errors"
	"os"
	"strconv"

	"github.com/rs/zerolog"
)

var (
	errIsPrettyNotFound = errors.New("is pretty not found")
	errVersionNotFound  = errors.New("version not found")
	errLogLevelNotFound = errors.New("log level not found")
)

const (
	isPrettyEnvName = "IS_PRETTY"
	versionEnvName  = "VERSION"
	logLevelEnvName = "LOG_LEVEL"
)

type LoggerConfig interface {
	IsPretty() bool
	Version() string
	LogLevel() zerolog.Level
}

type loggerConfig struct {
	isPretty bool
	version  string
	logLevel zerolog.Level
}

// NewLoggerConfig инициализирует конфиг логгера.
func NewLoggerConfig() (LoggerConfig, error) {
	isPretty, err := strconv.ParseBool(os.Getenv(isPrettyEnvName))
	if err != nil {
		return nil, errIsPrettyNotFound
	}

	version := os.Getenv(versionEnvName)
	if len(version) == 0 {
		return nil, errVersionNotFound
	}

	logLevel, err := zerolog.ParseLevel(os.Getenv(logLevelEnvName))
	if err != nil {
		return nil, errLogLevelNotFound
	}

	return &loggerConfig{
		isPretty: isPretty,
		version:  version,
		logLevel: logLevel,
	}, nil
}

func (cfg *loggerConfig) IsPretty() bool {
	return cfg.isPretty
}

func (cfg *loggerConfig) Version() string {
	return cfg.version
}

func (cfg *loggerConfig) LogLevel() zerolog.Level {
	return cfg.logLevel
}
