package user

import (
	"context"
	"errors"
	"user-service/internal/errorz"
	pb "user-service/pkg/api/user_v1"

	"github.com/rs/zerolog/log"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func (s *Server) GetUser(ctx context.Context, req *pb.GetUserRequest) (*pb.GetUserResponse, error) {
	// Логика получения пользователя
	user, err := s.userService.GetByID(ctx, req.GetId())
	switch {
	case errors.Is(err, errorz.ErrUserNotFound):
		return nil, status.Error(codes.NotFound, "User with this ID not found.")

	case err != nil:
		log.Error().Err(err).Msg("error while getting user")
		return nil, status.Error(codes.Internal, "Failed to get user.")
	}

	return &pb.GetUserResponse{
		User: &pb.User{
			Id:        user.ID,
			Email:     user.Email,
			Grade:     user.Grade,
			CreatedAt: timestamppb.New(user.CreatedAt),
		},
	}, nil
}
