package user

import (
	"context"
	"errors"
	"user-service/internal/errorz"
	"user-service/internal/models"

	"golang.org/x/crypto/bcrypt"
)

func (s *Service) Register(ctx context.Context, user *models.User) (string, error) {
	dbUser, err := s.userRepository.GetByEmail(ctx, user.Email)
	switch {
	case err == nil:
		return dbUser.ID, nil

	case !errors.Is(err, errorz.ErrUserNotFound):
		return "", err
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	user.Password = string(hashedPassword)

	return s.userRepository.Create(ctx, user)
}
