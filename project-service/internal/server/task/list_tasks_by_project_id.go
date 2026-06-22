package task

import (
	"context"
	"errors"
	"project-service/internal/errorz"
	pb "project-service/pkg/api/task_v1"

	"github.com/rs/zerolog/log"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func (s *Server) ListTasksByProjectID(ctx context.Context, req *pb.ListTasksByProjectIDRequest) (*pb.ListTasksByProjectIDResponse, error) {
	tasks, err := s.taskService.GetAllByProjectID(ctx, req.ProjectId)
	switch {
	case errors.Is(err, errorz.ErrUserIDNotSet):
		return nil, status.Error(codes.Unauthenticated, "Authentication required. Please provide a valid token.")

	case errors.Is(err, errorz.ErrProjectNotFound):
		return nil, status.Error(codes.NotFound, "Project with this ID not found.")

	case errors.Is(err, errorz.ErrTaskAccessForbidden):
		return nil, status.Error(codes.PermissionDenied, "Access denied.")

	case err != nil:
		log.Error().Err(err).Msg("error while listing tasks by project id")
		return nil, status.Errorf(codes.Internal, "Failed to list tasks.")
	}

	var response []*pb.Task
	for _, task := range tasks {
		response = append(response, &pb.Task{
			Id:          task.ID,
			Title:       task.Title,
			Description: task.Description,
			Status:      string(task.Status),
			ProjectId:   task.ProjectId,
			ExecutorId:  task.ExecutorId,
			Deadline:    timestamppb.New(task.Deadline),
			CreatedAt:   timestamppb.New(task.CreatedAt),
			UpdatedAt:   timestamppb.New(task.UpdatedAt),
		})
	}

	return &pb.ListTasksByProjectIDResponse{
		Tasks: response,
	}, nil
}
