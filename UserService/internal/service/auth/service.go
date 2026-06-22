package auth

import (
	"user-service/internal/config"
	"user-service/internal/repository"
)

type Service struct {
	jwtConfig config.JWTConfig

	userRepository repository.UserRepository
}

func NewService(userRepository repository.UserRepository, jwtConfig config.JWTConfig) *Service {
	return &Service{userRepository: userRepository, jwtConfig: jwtConfig}
}
