package integration

import (
	"context"
	"flag"
	"fmt"
	"testing"
	"time"

	"chat-service/test/integration/grpc"
	mongodb "chat-service/test/integration/mongo"
	rdb "chat-service/test/integration/redis"

	"github.com/ilyakaznacheev/cleanenv"
	"github.com/redis/go-redis/v9"
	"github.com/stretchr/testify/suite"
	"github.com/testcontainers/testcontainers-go/modules/compose"
	"github.com/testcontainers/testcontainers-go/wait"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const (
	Timeout = 5 * time.Minute
)

type Config struct {
	RedisConfig rdb.Config
	GrpcConfig  grpc.Config
	MongoConfig mongodb.Config
}

type BaseTestSuite struct {
	suite.Suite
	compose compose.ComposeStack
	mongo   *mongo.Client
	redis   *redis.Client
	config  *Config
	ctx     context.Context
	cancel  context.CancelFunc
}

func TestBase(t *testing.T) {
	suite.Run(t, new(BaseTestSuite))
}

func (s *BaseTestSuite) TestUser() {
	t := &UserSuite{base: s}
	suite.Run(s.T(), t)
}

func (s *BaseTestSuite) TestChat() {
	t := &ChatSuite{base: s}
	suite.Run(s.T(), t)
}

func (s *BaseTestSuite) TestMessage() {
	t := &MessageSuite{base: s}
	suite.Run(s.T(), t)
}

func (s *BaseTestSuite) SetupSuite() {
	var envPath, composePath string
	flag.StringVar(&envPath, "env", "../../test.env", "Path to test env file")
	flag.StringVar(&composePath, "compose", "../../docker-compose.test.yaml", "Path to docker compose file")
	flag.Parse()

	ctx, cancel := context.WithTimeout(context.Background(), Timeout)
	s.ctx = ctx
	s.cancel = cancel

	s.MustInitConfig(envPath)
	s.MustInitCompose(composePath)
	s.compose.Up(ctx)

	s.MustInitMongo()
	s.MustInitRedis()
}

func (s *BaseTestSuite) MustInitConfig(envPath string) {
	var config Config
	if err := cleanenv.ReadConfig(envPath, &config); err != nil {
		s.FailNow("read config: " + err.Error())
	}
	s.config = &config
}

func (s *BaseTestSuite) MustInitCompose(composePath string) {
	c, err := compose.NewDockerCompose(composePath)
	if err != nil {
		s.FailNow("init compose: " + err.Error())
	}
	s.compose = c.
		WaitForService("mongo-chat-test", wait.NewHealthStrategy()).
		WaitForService("redis-chat-test", wait.NewHealthStrategy()).
		WaitForService("chat-test", wait.NewLogStrategy("Start reading from stream").WithOccurrence(1))
}

func (s *BaseTestSuite) MustInitMongo() {
	uri := fmt.Sprintf("mongodb://%s:%s", s.config.MongoConfig.Host, s.config.MongoConfig.Port)

	var client *mongo.Client
	var err error
	for i := 0; i < 30; i++ {
		client, err = mongo.Connect(s.ctx, options.Client().ApplyURI(uri))
		if err == nil {
			err = client.Ping(s.ctx, nil)
			if err == nil {
				break
			}
		}
		time.Sleep(time.Second)
	}
	if err != nil {
		s.FailNow("init mongo: " + err.Error())
	}
	s.mongo = client
}

func (s *BaseTestSuite) MustInitRedis() {
	addr := fmt.Sprintf("%s:%s", s.config.RedisConfig.Host, s.config.RedisConfig.Port)
	client := redis.NewClient(&redis.Options{
		Addr:     addr,
		Password: "",
		DB:       0,
	})

	var err error
	for i := 0; i < 30; i++ {
		err = client.Ping(s.ctx).Err()
		if err == nil {
			break
		}
		time.Sleep(time.Second)
	}
	if err != nil {
		s.FailNow("init redis: " + err.Error())
	}
	s.redis = client
}

func (s *BaseTestSuite) SetupTest() {
	app, err := s.compose.ServiceContainer(s.ctx, "chat-test")
	if err != nil {
		s.FailNow("get container: " + err.Error())
	}
	err = app.Stop(s.ctx, nil)
	if err != nil {
		s.FailNow("stop container: " + err.Error())
	}

	if err := s.mongo.Database("chat").Drop(s.ctx); err != nil {
		s.FailNow("clean mongo: " + err.Error())
	}
	if err := s.redis.FlushAll(s.ctx).Err(); err != nil {
		s.FailNow("clean redis: " + err.Error())
	}

	err = app.Start(s.ctx)
	if err != nil {
		s.FailNow("start container: " + err.Error())
	}

	err = wait.ForListeningPort("50051/tcp").
		WithStartupTimeout(30 * time.Second).
		WaitUntilReady(s.ctx, app)
	if err != nil {
		s.FailNow("wait container: " + err.Error())
	}
}

func (s *BaseTestSuite) TearDownSuite() {
	s.compose.Down(s.ctx)
	s.mongo.Disconnect(s.ctx)
	s.redis.Close()
	s.cancel()
}
