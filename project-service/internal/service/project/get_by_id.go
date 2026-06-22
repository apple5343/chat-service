package project

import (
	"context"
	"project-service/internal/errorz"
	"project-service/internal/models"

	"google.golang.org/grpc/metadata"
)

func (s *Service) GetByID(ctx context.Context, id string) (*models.Project, error) {
	var userID string
	md, ok := metadata.FromIncomingContext(ctx)
	if ok {
		userID = md["user_id"][0]
	}
	if !ok {
		return nil, errorz.ErrUserIDNotSet
	}

	project, err := s.projectRepository.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	for _, pUserID := range project.Users {
		if userID == pUserID {
			return project, nil
		}
	}

	return nil, errorz.ErrProjectAccessForbidden
}
