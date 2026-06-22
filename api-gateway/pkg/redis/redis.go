package redis

import (
	"api-gateway/internal/config"

	"github.com/redis/go-redis/v9"
)

type Config struct {
	Addr     string `env:"REDIS_ADDR"`
	Password string `env:"REDIS_PASSWORD" default:""`
	DB       int    `env:"REDIS_DB" default:"0"`
}

func New(cfg config.RedisConfig) *redis.Client {
	return redis.NewClient(&redis.Options{
		Addr:     cfg.Address(),
		Password: cfg.Password(),
		DB:       cfg.Db(),
	})
}
