package main

import (
	"chat-service/internal/repository"
	"chat-service/internal/service"
	"chat-service/internal/transport/grpc"
	redisconsumer "chat-service/internal/transport/redis"
	"chat-service/pkg/logger"
	"chat-service/pkg/mongodb"
	"chat-service/pkg/redis"
	"context"

	"go.uber.org/zap"
	"golang.org/x/sync/errgroup"
)

func main() {
	log := logger.NewZapLogger("chat-service")
	context := logger.ContextWithLogger(context.Background(), log)

	redisCfg := redis.LoadConfig()
	rdb := redis.New(redisCfg)

	mongodbCfg := mongodb.LoadConfig()
	client := mongodb.New(mongodbCfg)

	updatesTransport := redisconsumer.NewUpdatesTransport(rdb)
	chatRepo := repository.NewChatRepository(&client)
	chatService := service.NewChatService(chatRepo, updatesTransport)

	redisConsumerCfg := redisconsumer.LoadConfig()
	redisConsumerServer := redisconsumer.New(redisConsumerCfg, rdb, chatService)

	grpcServer := grpc.New(grpc.LoadConfig(), chatService, context)

	eg := errgroup.Group{}

	eg.Go(func() error {
		return grpcServer.Start()
	})

	eg.Go(func() error {
		return redisConsumerServer.Run()
	})

	if err := eg.Wait(); err != nil {
		log.Error(context, "Error occurred", zap.Error(err))
	}
}
