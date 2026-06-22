package task

import (
	"context"
	"project-service/internal/errorz"
	"project-service/internal/models"

	"google.golang.org/grpc/metadata"
)

func (s *Service) GetAllByProjectID(ctx context.Context, projectID string) ([]*models.Task, error) {
	var userID string
	md, ok := metadata.FromIncomingContext(ctx)
	if ok {
		userID = md["user_id"][0]
	}
	if !ok {
		return nil, errorz.ErrUserIDNotSet
	}

	project, err := s.projectRepository.GetByID(ctx, projectID)
	if err != nil {
		return nil, err
	}

	for _, pUserID := range project.Users {
		if userID == pUserID {
			return s.taskRepository.GetAllByProjectID(ctx, projectID)
		}
	}

	return nil, errorz.ErrTaskAccessForbidden
}
