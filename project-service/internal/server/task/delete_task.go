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

func (s *Server) DeleteTask(ctx context.Context, req *pb.DeleteTaskRequest) (*pb.DeleteTaskResponse, error) {
	project, err := s.taskService.Delete(ctx, req.Id)
	switch {
	case errors.Is(err, errorz.ErrUserIDNotSet):
		return nil, status.Error(codes.Unauthenticated, "Authentication required. Please provide a valid token.")

	case errors.Is(err, errorz.ErrTaskNotFound):
		return nil, status.Error(codes.NotFound, "Task with this ID not found.")

	case errors.Is(err, errorz.ErrTaskAccessForbidden):
		return nil, status.Error(codes.PermissionDenied, "Access denied for this task.")

	case err != nil:
		log.Error().Err(err).Msg("error while deleting task")
		return nil, status.Errorf(codes.Internal, "Failed to delete task.")
	}

	for _, tgSubID := range project.NotificationSubscribersTGIDS {
		msg := fmt.Sprintf(models.DeleteTaskMsg, tgSubID, project.Title)
		if err = s.producer.SendMessage(s.kafkaConfig.Topic(), msg); err != nil {
			log.Error().Err(err).Msg("Failed to send message to Kafka topic")
		}
	}

	return &pb.DeleteTaskResponse{Success: true}, nil
}
