package project

import (
	"context"
	"project-service/internal/models"
)

func (s *Service) Create(ctx context.Context, project *models.Project) (string, error) {
	return s.projectRepository.Create(ctx, project)
}
