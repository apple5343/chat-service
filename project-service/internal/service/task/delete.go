package task

import (
	"context"
	"project-service/internal/errorz"
	"project-service/internal/models"

	"google.golang.org/grpc/metadata"
)

func (s *Service) Delete(ctx context.Context, id string) (*models.Project, error) {
	var userID string
	md, ok := metadata.FromIncomingContext(ctx)
	if ok {
		userID = md["user_id"][0]
	}
	if !ok {
		return nil, errorz.ErrUserIDNotSet
	}

	task, err := s.taskRepository.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	project, err := s.projectRepository.GetByID(ctx, task.ProjectId)
	if err != nil {
		return nil, err
	}

	if userID != project.AdminID {
		return nil, errorz.ErrTaskAccessForbidden
	}

	return project, s.taskRepository.Delete(ctx, id)
}
