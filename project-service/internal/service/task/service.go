package task

import "project-service/internal/repository"

type Service struct {
	taskRepository    repository.TaskRepository
	projectRepository repository.ProjectRepository
}

func NewService(taskRepository repository.TaskRepository, projectRepository repository.ProjectRepository) *Service {
	return &Service{taskRepository: taskRepository, projectRepository: projectRepository}
}
