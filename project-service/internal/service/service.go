package service

import (
	"context"
	"project-service/internal/models"
	projectpb "project-service/pkg/api/project_v1"
	taskpb "project-service/pkg/api/task_v1"
)

type ProjectService interface {
	Create(ctx context.Context, project *models.Project) (string, error)
	GetByID(ctx context.Context, id string) (*models.Project, error)
	GetAllByUserID(ctx context.Context, userID string) ([]*models.Project, error)
	Update(ctx context.Context, req *projectpb.UpdateProjectRequest) error
	AddUser(ctx context.Context, newUserID, projectID string) error
	Delete(ctx context.Context, id string) error
	SubscribeOnNotifications(ctx context.Context, projectID string, telegramID int64) error
}

type TaskService interface {
	Create(ctx context.Context, task *models.Task) (*models.Project, string, error)
	GetByID(ctx context.Context, id string) (*models.Task, error)
	GetAllByProjectID(ctx context.Context, projectID string) ([]*models.Task, error)
	Update(ctx context.Context, req *taskpb.UpdateTaskRequest) (*models.Project, error)
	Delete(ctx context.Context, id string) (*models.Project, error)
}
