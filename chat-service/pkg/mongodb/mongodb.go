package mongodb

import (
	"context"

	"github.com/ilyakaznacheev/cleanenv"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type Config struct {
	URI string `env:"MONGODB_URI"`
}

func LoadConfig() Config {
	var cfg Config
	err := cleanenv.ReadEnv(&cfg)
	if err != nil {
		panic(err)
	}
	return cfg
}

func New(cfg Config) mongo.Client {
	clientOpt := options.Client().ApplyURI(cfg.URI)
	client, err := mongo.Connect(context.Background(), clientOpt)
	if err != nil {
		panic(err)
	}
	err = client.Ping(context.Background(), nil)
	if err != nil {
		panic(err)
	}
	return *client
}
