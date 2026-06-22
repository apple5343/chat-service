package project

import "project-service/internal/repository"

type Service struct {
	projectRepository repository.ProjectRepository
}

func NewService(projectRepository repository.ProjectRepository) *Service {
	return &Service{projectRepository: projectRepository}
}
