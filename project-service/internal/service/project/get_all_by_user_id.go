package project

import (
	"context"
	"project-service/internal/models"
)

func (s *Service) GetAllByUserID(ctx context.Context, userID string) ([]*models.Project, error) {
	return s.projectRepository.GetAllByUserID(ctx, userID)
}
