package task

import (
	"context"
	"project-service/internal/errorz"
	"project-service/internal/models"
	pb "project-service/pkg/api/task_v1"

	"google.golang.org/grpc/metadata"
)

func (s *Service) Update(ctx context.Context, req *pb.UpdateTaskRequest) (*models.Project, error) {
	var userID string
	md, ok := metadata.FromIncomingContext(ctx)
	if ok {
		userID = md["user_id"][0]
	}
	if !ok {
		return nil, errorz.ErrUserIDNotSet
	}

	task, err := s.taskRepository.GetByID(ctx, req.GetId())
	if err != nil {
		return nil, err
	}

	project, err := s.projectRepository.GetByID(ctx, task.ProjectId)
	if err != nil {
		return nil, err
	}

	if userID != project.AdminID {
		return nil, errorz.ErrProjectAccessForbidden
	}

	s.update(task, req)

	return project, s.taskRepository.Update(ctx, task)
}

func (s *Service) update(task *models.Task, req *pb.UpdateTaskRequest) {
	if req.GetTitle() != "" {
		task.Title = req.GetTitle()
	}
	if req.GetDescription() != "" {
		task.Description = req.GetDescription()
	}
	if req.GetStatus() != "" {
		task.Status = models.TaskStatus(req.GetStatus())
	}
	if req.GetExecutor() != "" {
		task.ExecutorId = req.GetExecutor()
	}
	if req.GetDeadline() != nil {
		task.Deadline = req.GetDeadline().AsTime()
	}
}
