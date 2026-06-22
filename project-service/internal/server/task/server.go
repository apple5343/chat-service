package task

import (
	"project-service/internal/config"
	"project-service/internal/kafka"
	"project-service/internal/service"
	pb "project-service/pkg/api/task_v1"
)

type Server struct {
	pb.UnimplementedTaskServiceServer
	taskService service.TaskService

	producer kafka.Producer

	kafkaConfig config.KafkaConfig
}

func NewServer(taskService service.TaskService, producer kafka.Producer, kafkaConfig config.KafkaConfig) *Server {
	return &Server{
		taskService: taskService,
		producer:    producer,
		kafkaConfig: kafkaConfig,
	}
}
