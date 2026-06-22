package auth

import (
	"context"
	"user-service/internal/models"
	"user-service/internal/utils"
)

func (s *Service) GetAccessToken(ctx context.Context, refreshToken string) (string, error) {
	claims, err := utils.VerifyToken(refreshToken, []byte(s.jwtConfig.RefreshTokenSecret()))
	if err != nil {
		return "", err
	}

	accessToken, err := utils.GenerateToken(models.User{
		ID: claims.UserID,
	},
		[]byte(s.jwtConfig.AccessTokenSecret()),
		s.jwtConfig.AccessTokenExpiration(),
	)
	if err != nil {
		return "", err
	}

	return accessToken, nil
}
