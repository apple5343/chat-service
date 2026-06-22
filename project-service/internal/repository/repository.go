package repository

import (
	"context"
	"project-service/internal/models"
)

type ProjectRepository interface {
	Create(ctx context.Context, project *models.Project) (string, error)
	GetByID(ctx context.Context, id string) (*models.Project, error)
	GetAllByUserID(ctx context.Context, userID string) ([]*models.Project, error)
	Update(ctx context.Context, project *models.Project) error
	Delete(ctx context.Context, id string) error
}

type TaskRepository interface {
	Create(ctx context.Context, task *models.Task) (string, error)
	GetAllByProjectID(ctx context.Context, projectID string) ([]*models.Task, error)
	GetByID(ctx context.Context, id string) (*models.Task, error)
	Update(ctx context.Context, task *models.Task) error
	Delete(ctx context.Context, id string) error
}
