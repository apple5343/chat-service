package project

import (
	"context"
	"errors"
	"project-service/internal/errorz"
	pb "project-service/pkg/api/project_v1"

	"github.com/rs/zerolog/log"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
)

func (s *Server) UpdateProject(ctx context.Context, req *pb.UpdateProjectRequest) (*emptypb.Empty, error) {
	err := s.projectService.Update(ctx, req)
	switch {
	case errors.Is(err, errorz.ErrUserIDNotSet):
		return nil, status.Error(codes.Unauthenticated, "Authentication required. Please provide a valid token.")

	case errors.Is(err, errorz.ErrProjectNotFound):
		return nil, status.Error(codes.NotFound, "Project with this ID not found.")

	case errors.Is(err, errorz.ErrProjectAccessForbidden):
		return nil, status.Error(codes.PermissionDenied, "Access denied for this project.")

	case err != nil:
		log.Error().Err(err).Msg("error while updating project")
		return nil, status.Error(codes.Internal, "Failed to update project.")
	}

	return &emptypb.Empty{}, nil
}
