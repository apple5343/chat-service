package user

import "user-service/internal/repository"

type Service struct {
	userRepository repository.UserRepository
}

func NewService(userRepository repository.UserRepository) *Service {
	return &Service{userRepository: userRepository}
}
