package app

import (
	"context"
	"fmt"
	"project-service/internal/kafka"
	"project-service/internal/kafka/producer"
	"project-service/internal/repository"
	"project-service/internal/service"

	"github.com/jmoiron/sqlx"

	projectRepository "project-service/internal/repository/postgres"
	taskRepository "project-service/internal/repository/postgres"

	projectservice "project-service/internal/service/project"
	taskservice "project-service/internal/service/task"

	projectServer "project-service/internal/server/project"
	taskServer "project-service/internal/server/task"

	db "project-service/pkg/pg"

	"project-service/internal/closer"
	"project-service/internal/config"
)

type serviceProvider struct {
	pgConfig     config.PGConfig
	grpcConfig   config.GRPCConfig
	loggerConfig config.LoggerConfig
	kafkaConfig  config.KafkaConfig

	producer kafka.Producer

	db                *sqlx.DB
	projectRepository repository.ProjectRepository
	taskRepository    repository.TaskRepository

	projectService service.ProjectService
	taskService    service.TaskService

	projectServer *projectServer.Server
	taskServer    *taskServer.Server
}

func newServiceProvider() *serviceProvider {
	return &serviceProvider{}
}

func (s *serviceProvider) PGConfig() config.PGConfig {
	if s.pgConfig == nil {
		cfg, err := config.NewPGConfig()
		if err != nil {
			panic(fmt.Errorf("failed to get pg config: %w", err))
		}

		s.pgConfig = cfg
	}

	return s.pgConfig
}

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

func (s *serviceProvider) KafkaConfig() config.KafkaConfig {
	if s.kafkaConfig == nil {
		cfg := config.NewKafkaConfig()
		s.kafkaConfig = cfg
	}

	return s.kafkaConfig
}

func (s *serviceProvider) DBClient(_ context.Context) *sqlx.DB {
	if s.db == nil {
		database := db.Init(s.PGConfig().PGConn())
		closer.Add(database.Close)
		s.db = database
	}

	return s.db
}

func (s *serviceProvider) Producer(_ context.Context) kafka.Producer {
	if s.producer == nil {
		p := producer.NewProducer(s.KafkaConfig().Brokers(), s.KafkaConfig().SaramaConfig())
		closer.Add(p.Close)
		s.producer = p
	}

	return s.producer
}

func (s *serviceProvider) TaskRepository(ctx context.Context) repository.TaskRepository {
	if s.taskRepository == nil {
		s.taskRepository = taskRepository.NewTaskRepository(s.DBClient(ctx))
	}
	return s.taskRepository
}

func (s *serviceProvider) ProjectRepository(ctx context.Context) repository.ProjectRepository {
	if s.projectRepository == nil {
		s.projectRepository = projectRepository.NewProjectRepository(s.DBClient(ctx))
	}

	return s.projectRepository
}

func (s *serviceProvider) TaskService(ctx context.Context) service.TaskService {
	if s.taskService == nil {
		s.taskService = taskservice.NewService(s.TaskRepository(ctx), s.ProjectRepository(ctx))
	}

	return s.taskService
}

func (s *serviceProvider) ProjectService(ctx context.Context) service.ProjectService {
	if s.projectService == nil {
		s.projectService = projectservice.NewService(s.ProjectRepository(ctx))
	}

	return s.projectService
}

func (s *serviceProvider) TaskServer(ctx context.Context) *taskServer.Server {
	if s.taskServer == nil {
		s.taskServer = taskServer.NewServer(s.TaskService(ctx), s.Producer(ctx), s.KafkaConfig())
	}

	return s.taskServer
}

func (s *serviceProvider) ProjectServer(ctx context.Context) *projectServer.Server {
	if s.projectServer == nil {
		s.projectServer = projectServer.NewServer(s.ProjectService(ctx))
	}

	return s.projectServer
}
