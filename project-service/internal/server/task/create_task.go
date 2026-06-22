package task

import (
	"context"
	"errors"
	"fmt"
	"project-service/internal/errorz"
	"project-service/internal/models"
	pb "project-service/pkg/api/task_v1"

	"github.com/rs/zerolog/log"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (s *Server) CreateTask(ctx context.Context, req *pb.CreateTaskRequest) (*pb.CreateTaskResponse, error) {
	task := &models.Task{
		Title:       req.Title,
		Description: req.Description,
		Status:      models.TaskStatus(req.Status),
		ProjectId:   req.ProjectId,
		ExecutorId:  req.Executor,
		Deadline:    req.Deadline.AsTime(),
	}

	project, taskID, err := s.taskService.Create(ctx, task)
	switch {
	case errors.Is(err, errorz.ErrUserIDNotSet):
		return nil, status.Error(codes.Unauthenticated, "Authentication required. Please provide a valid token.")

	case errors.Is(err, errorz.ErrProjectNotFound):
		return nil, status.Error(codes.NotFound, "Project with this ID not found.")

	case errors.Is(err, errorz.ErrTaskAccessForbidden):
		return nil, status.Error(codes.PermissionDenied, "Access denied for this task.")

	case err != nil:
		log.Error().Err(err).Msg("error while creating task")
		return nil, status.Errorf(codes.Internal, "Failed to create task.")
	}

	for _, tgSubID := range project.NotificationSubscribersTGIDS {
		msg := fmt.Sprintf(models.CreateTaskMsg, tgSubID, project.Title, taskID)
		if err = s.producer.SendMessage(s.kafkaConfig.Topic(), msg); err != nil {
			log.Error().Err(err).Msg("Failed to send message to Kafka topic")
		}
	}

	return &pb.CreateTaskResponse{Id: taskID}, nil
}
