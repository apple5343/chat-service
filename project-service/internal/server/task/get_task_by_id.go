package task

import (
	"context"
	"errors"
	"project-service/internal/errorz"

	"github.com/rs/zerolog/log"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"

	pb "project-service/pkg/api/task_v1"
)

func (s *Server) GetTaskByID(ctx context.Context, req *pb.GetTaskByIDRequest) (*pb.GetTaskByIDResponse, error) {
	task, err := s.taskService.GetByID(ctx, req.Id)
	switch {
	case errors.Is(err, errorz.ErrUserIDNotSet):
		return nil, status.Error(codes.Unauthenticated, "Authentication required. Please provide a valid token.")

	case errors.Is(err, errorz.ErrTaskNotFound):
		return nil, status.Error(codes.NotFound, "Task with this ID not found.")

	case errors.Is(err, errorz.ErrTaskAccessForbidden):
		return nil, status.Error(codes.PermissionDenied, "Access denied for this task.")

	case err != nil:
		log.Error().Err(err).Msg("error while getting task by id")
		return nil, status.Errorf(codes.Internal, "Failed to get task.")
	}

	return &pb.GetTaskByIDResponse{
		Task: &pb.Task{
			Id:          task.ID,
			Title:       task.Title,
			Description: task.Description,
			Status:      string(task.Status),
			ProjectId:   task.ProjectId,
			ExecutorId:  task.ExecutorId,
			Deadline:    timestamppb.New(task.Deadline),
			CreatedAt:   timestamppb.New(task.CreatedAt),
			UpdatedAt:   timestamppb.New(task.UpdatedAt),
		},
	}, nil
}
