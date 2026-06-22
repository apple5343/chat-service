package auth

import (
	"user-service/internal/service"
	pb "user-service/pkg/api/auth_v1"
)

type Server struct {
	pb.UnimplementedAuthV1Server
	authService service.AuthService
}

func NewServer(authService service.AuthService) *Server {
	return &Server{
		authService: authService,
	}
}
