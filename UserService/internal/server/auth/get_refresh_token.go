package auth

import (
	"context"
	pb "user-service/pkg/api/auth_v1"

	"github.com/rs/zerolog/log"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (s *Server) GetRefreshToken(ctx context.Context, req *pb.GetRefreshTokenRequest) (*pb.GetRefreshTokenResponse, error) {
	refreshToken, err := s.authService.GetRefreshToken(ctx, req.GetRefreshToken())
	if err != nil {
		log.Error().Err(err).Msg("Error while getting refresh token")
		return nil, status.Error(codes.Unauthenticated, "Failed to get refresh token.")
	}

	return &pb.GetRefreshTokenResponse{RefreshToken: refreshToken}, nil
}
