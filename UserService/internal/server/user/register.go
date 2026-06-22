package user

import (
	"context"
	"user-service/internal/models"
	pb "user-service/pkg/api/user_v1"

	"github.com/rs/zerolog/log"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (s *Server) Register(ctx context.Context, req *pb.RegisterRequest) (*pb.RegisterResponse, error) {
	// Логика регистрации пользователя
	user := &models.User{
		Email:    req.GetEmail(),
		Password: req.GetPassword(),
		Grade:    req.GetGrade(),
	}

	userID, err := s.userService.Register(ctx, user)
	if err != nil {
		log.Error().Err(err).Msg("error while registering user")
		return nil, status.Error(codes.Internal, "Failed to register user.")
	}

	return &pb.RegisterResponse{
		Id: userID,
	}, nil
}
