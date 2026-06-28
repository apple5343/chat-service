package mongo

type Config struct {
	Host string `env:"MONGO_HOST" env-required:"true"`
	Port string `env:"MONGO_PORT" env-required:"true"`
}
