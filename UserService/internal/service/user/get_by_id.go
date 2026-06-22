package user

import (
	"context"
	"user-service/internal/errorz"
	"user-service/internal/models"

	"google.golang.org/grpc/metadata"
)

func (s *Service) GetByID(ctx context.Context, id string) (*models.User, error) {
	var userID string
	md, ok := metadata.FromIncomingContext(ctx)
	if ok {
		userID = md["user_id"][0]
	}
	if !ok {
		return nil, errorz.ErrUserIDNotSet
	}

	if userID != id {
		return nil, errorz.ErrUserAccessDenied
	}

	return s.userRepository.GetByID(ctx, id)
}
