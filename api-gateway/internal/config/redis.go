package config

import (
	"errors"
	"net"
	"os"
	"strconv"
)

const (
	redisHostEnvName = "REDIS_HOST"
	redisPortEnvName = "REDIS_PORT"
	redisPasswordEnv = "REDIS_PASSWORD"
	redisDBEnvName   = "REDIS_DB"
)

var (
	errRedisHostNotFound = errors.New("redis host not found")
	errRedisPortNotFound = errors.New("redis port not found")
)

type RedisConfig interface {
	Address() string
	Password() string
	Db() int
}

type redisConfig struct {
	host     string
	port     string
	password string
	db       int
}

func (cfg *redisConfig) Address() string {
	return net.JoinHostPort(cfg.host, cfg.port)
}

func (cfg *redisConfig) Password() string {
	return cfg.password
}

func (cfg *redisConfig) Db() int {
	return cfg.db
}

// NewRedisConfig инициализирует конфиг редиса.
func NewRedisConfig() (RedisConfig, error) {
	host := os.Getenv(redisHostEnvName)
	if len(host) == 0 {
		return nil, errRedisHostNotFound
	}

	port := os.Getenv(redisPortEnvName)
	if len(port) == 0 {
		return nil, errRedisPortNotFound
	}

	password := os.Getenv(redisPasswordEnv)

	db := os.Getenv(redisDBEnvName)
	dbInt, _ := strconv.Atoi(db)

	return &redisConfig{
		host:     host,
		port:     port,
		password: password,
		db:       dbInt,
	}, nil
}
