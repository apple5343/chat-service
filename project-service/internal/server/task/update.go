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

func (s *Server) UpdateTask(ctx context.Context, req *pb.UpdateTaskRequest) (*pb.UpdateTaskResponse, error) {
	project, err := s.taskService.Update(ctx, req)
	switch {
	case errors.Is(err, errorz.ErrUserIDNotSet):
		return nil, status.Error(codes.Unauthenticated, "Authentication required. Please provide a valid token.")

	case errors.Is(err, errorz.ErrTaskNotFound):
		return nil, status.Error(codes.NotFound, "Task with this ID not found.")

	case errors.Is(err, errorz.ErrProjectAccessForbidden):
		return nil, status.Error(codes.PermissionDenied, "Access denied for this task.")

	case err != nil:
		log.Error().Err(err).Msg("error while updating task")
		return nil, status.Error(codes.Internal, "Failed to update task.")
	}

	for _, tgSubID := range project.NotificationSubscribersTGIDS {
		msg := fmt.Sprintf(models.UpdateTaskMsg, tgSubID, project.Title)
		if err = s.producer.SendMessage(s.kafkaConfig.Topic(), msg); err != nil {
			log.Error().Err(err).Msg("Failed to send message to Kafka topic")
		}
	}

	return &pb.UpdateTaskResponse{Success: true}, nil
}
