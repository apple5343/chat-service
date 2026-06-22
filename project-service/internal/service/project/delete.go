package project

import (
	"context"
	"project-service/internal/errorz"

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

	project, err := s.projectRepository.GetByID(ctx, id)
	if err != nil {
		return err
	}

	if userID != project.AdminID {
		return errorz.ErrProjectAccessForbidden
	}

	return s.projectRepository.Delete(ctx, id)
}
