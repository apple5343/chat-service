package project

import (
	"context"
	"project-service/internal/errorz"
	"project-service/internal/models"
	pb "project-service/pkg/api/project_v1"

	"google.golang.org/grpc/metadata"
)

func (s *Service) Update(ctx context.Context, req *pb.UpdateProjectRequest) error {
	var userID string
	md, ok := metadata.FromIncomingContext(ctx)
	if ok {
		userID = md["user_id"][0]
	}
	if !ok {
		return errorz.ErrUserIDNotSet
	}

	project, err := s.projectRepository.GetByID(ctx, req.GetId())
	if err != nil {
		return err
	}

	if userID != project.AdminID {
		return errorz.ErrProjectAccessForbidden
	}

	s.update(project, req)

	return s.projectRepository.Update(ctx, project)
}

func (s *Service) update(project *models.Project, req *pb.UpdateProjectRequest) {
	if req.GetTitle() != "" {
		project.Title = req.GetTitle()
	}
	if req.GetDescription() != "" {
		project.Description = req.GetDescription()
	}
	if req.GetAdminId() != "" {
		project.AdminID = req.GetAdminId()
	}
}
