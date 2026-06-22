package auth

import (
	"context"
	"user-service/internal/errorz"
	"user-service/internal/models"
	"user-service/internal/utils"
	pb "user-service/pkg/api/auth_v1"
)

func (s *Service) Login(ctx context.Context, req *pb.LoginRequest) (string, error) {
	user, err := s.userRepository.GetByEmail(ctx, req.GetEmail())
	if err != nil {
		return "", err
	}

	isMatch := utils.VerifyPassword(user.Password, req.GetPassword())
	if !isMatch {
		return "", errorz.ErrPasswordDoesNotMatch
	}

	refreshToken, err := utils.GenerateToken(
		models.User{ID: user.ID},
		[]byte(s.jwtConfig.RefreshTokenSecret()),
		s.jwtConfig.RefreshTokenExpiration(),
	)
	if err != nil {
		return "", err
	}

	return refreshToken, nil
}
