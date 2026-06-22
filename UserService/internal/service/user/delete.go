package user

import (
	"context"
	"user-service/internal/errorz"

	"google.golang.org/grpc/metadata"
)

func (s *Service) Delete(ctx context.Context, id string) error {
	var userID string
	md, ok := metadata.FromIncomingContext(ctx)
	if ok {
		userID = md["user_id"][0]
	}
	if !ok {
		return errorz.ErrUserIDNotSet
	}

	if userID != id {
		return errorz.ErrUserAccessDenied
	}

	return s.userRepository.Delete(ctx, id)
}
