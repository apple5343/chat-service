package service

import (
	"context"
	"user-service/internal/models"
	authpb "user-service/pkg/api/auth_v1"
	userpb "user-service/pkg/api/user_v1"
)

type UserService interface {
	Register(ctx context.Context, user *models.User) (string, error)
	GetByID(ctx context.Context, id string) (*models.User, error)
	Update(ctx context.Context, req *userpb.UpdateUserRequest) error
	Delete(ctx context.Context, id string) error
}

type AuthService interface {
	Login(ctx context.Context, req *authpb.LoginRequest) (string, error)
	GetRefreshToken(ctx context.Context, refreshToken string) (string, error)
	GetAccessToken(ctx context.Context, refreshToken string) (string, error)
}
