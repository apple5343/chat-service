package task

import (
	"context"
	"project-service/internal/errorz"
	"project-service/internal/models"

	"google.golang.org/grpc/metadata"
)

func (s *Service) Create(ctx context.Context, task *models.Task) (*models.Project, string, error) {
	var userID string
	md, ok := metadata.FromIncomingContext(ctx)
	if ok {
		userID = md["user_id"][0]
	}
	if !ok {
		return nil, "", errorz.ErrUserIDNotSet
	}

	project, err := s.projectRepository.GetByID(ctx, task.ProjectId)
	if err != nil {
		return nil, "", err
	}

	for _, pUserID := range project.Users {
		if userID == pUserID {
			var taskID string
			taskID, err = s.taskRepository.Create(ctx, task)
			return project, taskID, err
		}
	}

	return nil, "", errorz.ErrTaskAccessForbidden
}
