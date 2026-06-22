package app

import (
	"api-gateway/internal/closer"
	"api-gateway/internal/config"
	"api-gateway/internal/server/auth"
	"api-gateway/internal/server/chat"
	"api-gateway/internal/server/chat/ws"
	"api-gateway/internal/server/project"
	"api-gateway/internal/server/task"
	"api-gateway/internal/server/user"
	"api-gateway/internal/service"
	authpb "api-gateway/pkg/api/auth_v1"
	chatpb "api-gateway/pkg/api/chat_v1"
	projectpb "api-gateway/pkg/api/project_v1"
	taskpb "api-gateway/pkg/api/task_v1"
	userpb "api-gateway/pkg/api/user_v1"
	rdb "api-gateway/pkg/redis"
	"fmt"

	"github.com/redis/go-redis/v9"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

// Структура сервис провайдера.
type serviceProvider struct {
	grpcConfig   config.GRPCConfig
	httpConfig   config.HTTPConfig
	loggerConfig config.LoggerConfig
	jwtConfig    config.JWTConfig
	redisConfig  config.RedisConfig

	redis *redis.Client

	chatService *service.Chat
	hub         *ws.Hub

	userClientConn *grpc.ClientConn
	authAPIClient  authpb.AuthV1Client
	userAPIClient  userpb.UserV1Client

	projectClientConn *grpc.ClientConn
	projectAPIClient  projectpb.ProjectServiceClient
	taskAPIClient     taskpb.TaskServiceClient
	chatAPIClient     chatpb.ChatServiceClient
	chatClientConn    *grpc.ClientConn

	userHandler    *user.Handler
	authHandler    *auth.Handler
	projectHandler *project.Handler
	taskHandler    *task.Handler
	chatHandler    *chat.Handler
}

// Инициализирует сервис провайдер.
func newServiceProvider() *serviceProvider {
	return &serviceProvider{}
}

func (s *serviceProvider) Hub() *ws.Hub {
	if s.hub == nil {
		s.hub = ws.NewHub(s.ChatService())
	}
	return s.hub
}

// GRPCConfig геттер для grpc конфига.
func (s *serviceProvider) GRPCConfig() config.GRPCConfig {
	if s.grpcConfig == nil {
		cfg, err := config.NewGRPCConfig()
		if err != nil {
			panic(fmt.Errorf("failed to get grpc config: %w", err))
		}

		s.grpcConfig = cfg
	}

	return s.grpcConfig
}

// HTTPConfig геттер для http конфига.
func (s *serviceProvider) HTTPConfig() config.HTTPConfig {
	if s.httpConfig == nil {
		cfg, err := config.NewHTTPConfig()
		if err != nil {
			panic(fmt.Errorf("failed to get http config: %w", err))
		}
		s.httpConfig = cfg
	}

	return s.httpConfig
}

// LoggerConfig геттер для конфига логгера.
func (s *serviceProvider) LoggerConfig() config.LoggerConfig {
	if s.loggerConfig == nil {
		cfg, err := config.NewLoggerConfig()
		if err != nil {
			panic(fmt.Errorf("failed to get logger config: %w", err))
		}

		s.loggerConfig = cfg
	}

	return s.loggerConfig
}

// JWTConfig геттер для jwt конфига.
func (s *serviceProvider) JWTConfig() config.JWTConfig {
	if s.jwtConfig == nil {
		cfg, err := config.NewJWTConfig()
		if err != nil {
			panic(fmt.Errorf("failed to get jwt config: %w", err))
		}

		s.jwtConfig = cfg
	}

	return s.jwtConfig
}

// RedisConfig геттер для redis конфига.
func (s *serviceProvider) RedisConfig() config.RedisConfig {
	if s.redisConfig == nil {
		cfg, err := config.NewRedisConfig()
		if err != nil {
			panic(fmt.Errorf("failed to get redis config: %w", err))
		}

		s.redisConfig = cfg
	}

	return s.redisConfig
}

// Redis геттер для redis.
func (s *serviceProvider) Redis() *redis.Client {
	if s.redis == nil {
		cfg := s.RedisConfig()
		s.redis = rdb.New(cfg)
		closer.Add(s.redis.Close)
	}
	return s.redis
}

// ChatService геттер для сервиса чата.
func (s *serviceProvider) ChatService() *service.Chat {
	if s.chatService == nil {
		s.chatService = service.NewChat(s.Redis())
	}
	return s.chatService
}

// UserClientConn геттер для connection к клиенту User Servic'а.
func (s *serviceProvider) UserClientConn() *grpc.ClientConn {
	if s.userClientConn == nil {
		conn, err := grpc.NewClient(s.GRPCConfig().UserServerAddress(), grpc.WithTransportCredentials(insecure.NewCredentials()))
		if err != nil {
			panic(fmt.Errorf("failed to connect to gRPC user server: %w", err))
		}
		closer.Add(conn.Close)
		s.userClientConn = conn
	}
	return s.userClientConn
}

// AuthAPIClient геттер для клиента auth api.
func (s *serviceProvider) AuthAPIClient() authpb.AuthV1Client {
	if s.authAPIClient == nil {
		s.authAPIClient = authpb.NewAuthV1Client(s.UserClientConn())
	}
	return s.authAPIClient
}

// UserAPIClient геттер для клиента user api.
func (s *serviceProvider) UserAPIClient() userpb.UserV1Client {
	if s.userAPIClient == nil {
		s.userAPIClient = userpb.NewUserV1Client(s.UserClientConn())
	}
	return s.userAPIClient
}

// ProjectClientConn геттер для connection к клиенту Project Servic'а.
func (s *serviceProvider) ProjectClientConn() *grpc.ClientConn {
	if s.projectClientConn == nil {
		conn, err := grpc.NewClient(s.GRPCConfig().ProjectServerAddress(), grpc.WithTransportCredentials(insecure.NewCredentials()))
		if err != nil {
			panic(fmt.Errorf("failed to connect to gRPC project server: %w", err))
		}
		closer.Add(conn.Close)
		s.projectClientConn = conn
	}
	return s.projectClientConn
}

// ChatClientConn геттер для connection к клиенту Chat Servic'а.
func (s *serviceProvider) ChatClientConn() *grpc.ClientConn {
	if s.chatClientConn == nil {
		conn, err := grpc.NewClient(s.GRPCConfig().ChatServerAddress(), grpc.WithTransportCredentials(insecure.NewCredentials()))
		if err != nil {
			panic(fmt.Errorf("failed to connect to gRPC chat server: %w", err))
		}
		closer.Add(conn.Close)
		s.chatClientConn = conn
	}
	return s.chatClientConn
}

// ProjectAPIClient геттер для клиента для project api.
func (s *serviceProvider) ProjectAPIClient() projectpb.ProjectServiceClient {
	if s.projectAPIClient == nil {
		s.projectAPIClient = projectpb.NewProjectServiceClient(s.ProjectClientConn())
	}
	return s.projectAPIClient
}

// TaskAPIClient геттер для клиента для task api.
func (s *serviceProvider) TaskAPIClient() taskpb.TaskServiceClient {
	if s.taskAPIClient == nil {
		s.taskAPIClient = taskpb.NewTaskServiceClient(s.ProjectClientConn())
	}
	return s.taskAPIClient
}

// ChatAPIClient геттер для клиента для chat api.
func (s *serviceProvider) ChatAPIClient() chatpb.ChatServiceClient {
	if s.chatAPIClient == nil {
		s.chatAPIClient = chatpb.NewChatServiceClient(s.ChatClientConn())
	}
	return s.chatAPIClient
}

// UserHandler геттер для User Handler'а.
func (s *serviceProvider) UserHandler() *user.Handler {
	if s.userHandler == nil {
		s.userHandler = user.NewHandler(s.UserAPIClient())
	}
	return s.userHandler
}

// AuthHandler геттер для Auth Handler'а.
func (s *serviceProvider) AuthHandler() *auth.Handler {
	if s.authHandler == nil {
		s.authHandler = auth.NewHandler(s.AuthAPIClient(), s.JWTConfig())
	}
	return s.authHandler
}

// ProjectHandler геттер для Project Handler'а.
func (s *serviceProvider) ProjectHandler() *project.Handler {
	if s.projectHandler == nil {
		s.projectHandler = project.NewHandler(s.ProjectAPIClient(), s.ChatAPIClient())
	}
	return s.projectHandler
}

// TaskHandler геттер для Task Handler'а.
func (s *serviceProvider) TaskHandler() *task.Handler {
	if s.taskHandler == nil {
		s.taskHandler = task.NewHandler(s.TaskAPIClient())
	}
	return s.taskHandler
}

// ChatHandler геттер для Chat Handler'а.
func (s *serviceProvider) ChatHandler() *chat.Handler {
	if s.chatHandler == nil {
		s.chatHandler = chat.NewHandler(s.ChatAPIClient(), s.ChatService(), s.Hub())
	}
	return s.chatHandler
}
