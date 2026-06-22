package user

import (
	"context"
	"user-service/internal/errorz"
	"user-service/internal/models"
	pb "user-service/pkg/api/user_v1"

	"golang.org/x/crypto/bcrypt"
	"google.golang.org/grpc/metadata"
)

func (s *Service) Update(ctx context.Context, req *pb.UpdateUserRequest) error {
	var userID string
	md, ok := metadata.FromIncomingContext(ctx)
	if ok {
		userID = md["user_id"][0]
	}
	if !ok {
		return errorz.ErrUserIDNotSet
	}

	if req.GetId() != userID {
		return errorz.ErrUserAccessDenied
	}

	user, err := s.userRepository.GetByID(ctx, userID)
	if err != nil {
		return err
	}

	if err = s.update(user, req); err != nil {
		return err
	}

	return s.userRepository.Update(ctx, user)
}

func (s *Service) update(user *models.User, req *pb.UpdateUserRequest) error {
	if req.GetPassword() != "" {
		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.GetPassword()), bcrypt.DefaultCost)
		if err != nil {
			return err
		}
		user.Password = string(hashedPassword)
	}
	if req.GetGrade() != "" {
		user.Grade = req.GetGrade()
	}

	return nil
}
