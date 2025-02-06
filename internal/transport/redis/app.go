package redisconsumer

import (
	"chat-service/entity"
	"context"
	"fmt"
	"log"
	"sync"

	"github.com/ilyakaznacheev/cleanenv"
	"github.com/redis/go-redis/v9"
)

type ChatService interface {
	WriteMessage(ctx context.Context, msg *entity.Message)
	ReadMessage(ch chan *entity.Message)
}

type Config struct {
	StreamName   string `env:"STREAM_NAME"`
	GroupName    string `env:"GROUP_NAME"`
	ConsumerName string `env:"CONSUMER_NAME"`
	PublishChan  string `env:"PUBLISH_CHANNEL"`
	WorkerCount  int    `env:"WORKER_COUNT" default:"1"`
}

func LoadConfig() Config {
	var cfg Config
	err := cleanenv.ReadEnv(&cfg)
	if err != nil {
		panic(err)
	}
	return cfg
}

func New(config Config, rdb *redis.Client, chatService ChatService) *Server {
	ctx, cancel := context.WithCancel(context.Background())

	return &Server{
		redis:       rdb,
		chatService: chatService,
		cfg:         config,
		cancel:      cancel,
		ctx:         ctx,
		wg:          sync.WaitGroup{},
	}
}

type Server struct {
	redis       *redis.Client
	chatService ChatService
	cfg         Config
	cancel      context.CancelFunc
	ctx         context.Context
	wg          sync.WaitGroup
}

func (s *Server) Run() error {
	if err := s.createConsumerGroup(s.ctx, s.redis); err != nil {
		return err
	}
	sendCh := make(chan *entity.Message, 10)
	s.chatService.ReadMessage(sendCh)
	go s.StartSendMessagesProcessor(sendCh)
	go func(ch chan *entity.Message) {
		<-s.ctx.Done()
		close(ch)
	}(sendCh)
	fmt.Println("Start reading from stream")
	for {
		messages, err := s.redis.XReadGroup(s.ctx, &redis.XReadGroupArgs{
			Group:    s.cfg.GroupName,
			Consumer: s.cfg.ConsumerName,
			Streams:  []string{s.cfg.StreamName, ">"},
			Count:    10,
			Block:    0,
		}).Result()
		if err != nil {
			if err == redis.Nil {
				continue
			}
			if s.ctx.Err() != nil {
				return nil
			}
			log.Printf("Error reading from stream: %v\n", err)
			continue
		}
		for _, stream := range messages {
			for _, message := range stream.Messages {
				go func(msg *redis.XMessage) {
					s.wg.Add(1)
					s.processMessage(*msg)
					s.wg.Done()
				}(&message)
				if err := s.redis.XAck(context.Background(), s.cfg.StreamName, s.cfg.GroupName, message.ID).Err(); err != nil {
					log.Printf("Error acking message: %v\n", err)
				}
			}
		}
	}
}

func (s *Server) Shutdown() {
	s.cancel()
	s.wg.Wait()
}

func (s *Server) createConsumerGroup(ctx context.Context, rdb *redis.Client) error {
	err := rdb.XGroupCreateMkStream(ctx, s.cfg.StreamName, s.cfg.GroupName, "$").Err()
	if err != nil && err.Error() != "BUSYGROUP Consumer Group name already exists" {
		return err
	}
	log.Println("Consumer Group ready.")
	return nil
}
