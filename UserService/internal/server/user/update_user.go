package user

import (
	"context"
	"errors"
	"user-service/internal/errorz"
	pb "user-service/pkg/api/user_v1"

	"github.com/rs/zerolog/log"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
)

func (s *Server) UpdateUser(ctx context.Context, req *pb.UpdateUserRequest) (*emptypb.Empty, error) {
	err := s.userService.Update(ctx, req)
	switch {
	case errors.Is(err, errorz.ErrUserIDNotSet):
		return nil, status.Error(codes.Unauthenticated, "User id not set.")
	case errors.Is(err, errorz.ErrUserAccessDenied):
		return nil, status.Error(codes.PermissionDenied, "Permission denied.")
	case errors.Is(err, errorz.ErrUserNotFound):
		return nil, status.Error(codes.NotFound, "User not found.")
	case err != nil:
		log.Error().Err(err).Msg("error while updating user")
		return nil, status.Error(codes.Internal, "Failed to update user")
	}

	return &emptypb.Empty{}, nil
}
