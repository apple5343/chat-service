package auth

import (
	"context"
	"user-service/internal/models"
	"user-service/internal/utils"
)

func (s *Service) GetRefreshToken(ctx context.Context, refreshToken string) (string, error) {
	claims, err := utils.VerifyToken(refreshToken, []byte(s.jwtConfig.RefreshTokenSecret()))
	if err != nil {
		return "", err
	}

	refreshToken, err = utils.GenerateToken(
		models.User{ID: claims.UserID},
		[]byte(s.jwtConfig.RefreshTokenSecret()),
		s.jwtConfig.RefreshTokenExpiration(),
	)
	if err != nil {
		return "", err
	}

	return refreshToken, nil
}
