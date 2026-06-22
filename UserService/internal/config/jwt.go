package config

import (
	"errors"
	"fmt"
	"os"
	"time"
)

var (
	errRefreshTokenSecretNotFound     = errors.New("refresh token secret not found")
	errRefreshTokenExpirationNotFound = errors.New("refresh token expiration not found")
	errAccessTokenSecretNotFound      = errors.New("access token secret not found")
	errAccessTokenExpirationNotFound  = errors.New("access token expiration not found")
)

const (
	refreshTokenSecretEnvName     = "REFRESH_TOKEN_SECRET"
	refreshTokenExpirationEnvName = "REFRESH_TOKEN_EXPIRATION"
	accessTokenSecretEnvName      = "ACCESS_TOKEN_SECRET"
	accessTokenExpirationEnvName  = "ACCESS_TOKEN_EXPIRATION"
)

type JWTConfig interface {
	RefreshTokenSecret() string
	RefreshTokenExpiration() time.Duration
	AccessTokenSecret() string
	AccessTokenExpiration() time.Duration
}

type jwtConfig struct {
	refreshTokenSecret     string
	refreshTokenExpiration time.Duration
	accessTokenSecret      string
	accessTokenExpiration  time.Duration
}

func NewJWTConfig() (JWTConfig, error) {
	refreshTokenSecret := os.Getenv(refreshTokenSecretEnvName)
	if len(refreshTokenSecret) == 0 {
		return nil, errRefreshTokenSecretNotFound
	}

	refreshTokenExpirationStr := os.Getenv(refreshTokenExpirationEnvName)
	if len(refreshTokenSecret) == 0 {
		return nil, errRefreshTokenExpirationNotFound
	}
	refreshTokenExpiration, err := time.ParseDuration(refreshTokenExpirationStr)
	if err != nil {
		return nil, fmt.Errorf("parse refresh token expiration: %w", err)
	}

	accessTokenSecret := os.Getenv(accessTokenSecretEnvName)
	if len(accessTokenSecret) == 0 {
		return nil, errAccessTokenSecretNotFound
	}

	accessTokenExpirationStr := os.Getenv(accessTokenExpirationEnvName)
	if len(accessTokenExpirationStr) == 0 {
		return nil, errAccessTokenExpirationNotFound

	}
	accessTokenExpiration, err := time.ParseDuration(accessTokenExpirationStr)
	if err != nil {
		return nil, fmt.Errorf("parse access token expiration: %w", err)
	}

	return &jwtConfig{
		refreshTokenSecret:     refreshTokenSecret,
		refreshTokenExpiration: refreshTokenExpiration,
		accessTokenSecret:      accessTokenSecret,
		accessTokenExpiration:  accessTokenExpiration,
	}, nil
}

func (cfg *jwtConfig) RefreshTokenSecret() string {
	return cfg.refreshTokenSecret
}

func (cfg *jwtConfig) RefreshTokenExpiration() time.Duration {
	return cfg.refreshTokenExpiration
}

func (cfg *jwtConfig) AccessTokenSecret() string {
	return cfg.accessTokenSecret
}

func (cfg *jwtConfig) AccessTokenExpiration() time.Duration {
	return cfg.accessTokenExpiration
}
