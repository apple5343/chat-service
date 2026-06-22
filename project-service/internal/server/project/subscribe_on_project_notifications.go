package project

import (
	"context"
	"errors"
	"project-service/internal/errorz"
	pb "project-service/pkg/api/project_v1"

	"github.com/rs/zerolog/log"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (s *Server) SubscribeOnProjectNotifications(ctx context.Context, req *pb.SubscribeOnProjectNotificationsRequest) (*pb.SubscribeOnProjectNotificationsResponse, error) {
	err := s.projectService.SubscribeOnNotifications(ctx, req.GetProjectID(), req.GetTelegramID())
	switch {
	case errors.Is(err, errorz.ErrUserIDNotSet):
		return nil, status.Error(codes.Unauthenticated, "Authentication required. Please provide a valid token.")

	case errors.Is(err, errorz.ErrProjectNotFound):
		return nil, status.Error(codes.NotFound, "Project with this ID not found.")

	case errors.Is(err, errorz.ErrProjectAccessForbidden):
		return nil, status.Error(codes.PermissionDenied, "Access denied for this project.")

	case err != nil:
		log.Error().Err(err).Msg("error while subscribing on project notifications")
		return nil, status.Error(codes.Internal, "Failed to subscribe on project notifications.")
	}

	return &pb.SubscribeOnProjectNotificationsResponse{
		Msg: "To receive notifications, launch the bot: @taskly_notifications_bot",
	}, nil
}
