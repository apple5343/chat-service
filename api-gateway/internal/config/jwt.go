package config

import (
	"errors"
	"os"
)

var (
	errAccessTokenSecretNotFound = errors.New("access token secret not found")
)

const (
	accessTokenSecretEnvName = "ACCESS_TOKEN_SECRET"
)

type JWTConfig interface {
	AccessTokenSecret() string
}

type jwtConfig struct {
	accessTokenSecret string
}

// NewJWTConfig инициализирует jwt конфиг.
func NewJWTConfig() (JWTConfig, error) {
	accessTokenSecret := os.Getenv(accessTokenSecretEnvName)
	if len(accessTokenSecret) == 0 {
		return nil, errAccessTokenSecretNotFound
	}

	return &jwtConfig{
		accessTokenSecret: accessTokenSecret,
	}, nil
}

func (cfg *jwtConfig) AccessTokenSecret() string {
	return cfg.accessTokenSecret
}
