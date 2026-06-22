package project

import (
	"project-service/internal/service"

	pb "project-service/pkg/api/project_v1"
)

type Server struct {
	pb.UnimplementedProjectServiceServer
	projectService service.ProjectService
}

func NewServer(projectService service.ProjectService) *Server {
	return &Server{
		projectService: projectService,
	}
}
