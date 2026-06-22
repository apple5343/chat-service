package config

import (
	"errors"
	"os"
	"strconv"

	"github.com/rs/zerolog"
)

var (
	ErrIsPrettyNotFound = errors.New("is pretty not found")
	ErrVersionNotFound  = errors.New("version not found")
	ErrLogLevelNotFound = errors.New("log level not found")
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

func NewLoggerConfig() (LoggerConfig, error) {
	isPretty, err := strconv.ParseBool(os.Getenv(isPrettyEnvName))
	if err != nil {
		return nil, ErrIsPrettyNotFound
	}

	version := os.Getenv(versionEnvName)
	if len(version) == 0 {
		return nil, ErrVersionNotFound
	}

	logLevel, err := zerolog.ParseLevel(os.Getenv(logLevelEnvName))
	if err != nil {
		return nil, ErrLogLevelNotFound
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
