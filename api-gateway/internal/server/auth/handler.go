package auth

import (
	"api-gateway/internal/config"
	authpb "api-gateway/pkg/api/auth_v1"
)

// Handler структура handler'а пользователей.
type Handler struct {
	authAPIClient authpb.AuthV1Client

	jwtConfig config.JWTConfig
}

// NewHandler инициализирует Handler.
func NewHandler(authAPIClient authpb.AuthV1Client, jwtConfig config.JWTConfig) *Handler {
	return &Handler{authAPIClient: authAPIClient, jwtConfig: jwtConfig}
}
