package grpc

import (
	api "chat-service/pkg/api/chat_v1"
	"chat-service/pkg/logger"
	"context"
	"fmt"
	"net"

	"github.com/ilyakaznacheev/cleanenv"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

type Server struct {
	api.UnsafeChatServiceServer
	grpcServer  *grpc.Server
	listener    net.Listener
	cfg         Config
	chatService ChatService
}

type Config struct {
	Port string `env:"GRPC_PORT"`
}

func LoadConfig() Config {
	var cfg Config
	err := cleanenv.ReadEnv(&cfg)
	if err != nil {
		panic(err)
	}
	return cfg
}

func New(cfg Config, chatService ChatService, ctx context.Context) *Server {
	lis, err := net.Listen("tcp", fmt.Sprintf(":%s", cfg.Port))
	if err != nil {
		panic(err)
	}

	opts := []grpc.ServerOption{
		grpc.UnaryInterceptor(RequestsLogger(logger.GetLogger(ctx))),
	}

	grpcServer := grpc.NewServer(opts...)
	s := &Server{
		cfg:         cfg,
		listener:    lis,
		grpcServer:  grpcServer,
		chatService: chatService,
	}
	reflection.Register(grpcServer)
	api.RegisterChatServiceServer(grpcServer, s)
	return s
}

func (s *Server) Start() error {
	return s.grpcServer.Serve(s.listener)
}

func (s *Server) Shutdown() {
	s.grpcServer.GracefulStop()
}
