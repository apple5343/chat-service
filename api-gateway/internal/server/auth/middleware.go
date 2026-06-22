package auth

import (
	"api-gateway/internal/errorz"
	"api-gateway/internal/server/helper"
	"api-gateway/internal/utils"
	"context"
	"fmt"
	"net/http"
	"strings"

	"google.golang.org/grpc/metadata"
)

// Auth middleware проверяет, авторизован ли пользователь.
func (h *Handler) Auth(next http.Handler) http.Handler {
	return helper.MakeHandler(func(w http.ResponseWriter, r *http.Request) error {
		// Получение контекста из запроса
		ctx := r.Context()

		// Извлечение заголовка Authorization из запроса
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			return errorz.APIError{
				Status: http.StatusUnauthorized,
				Err:    fmt.Errorf("empty authorization header"),
				Msg:    "auth header is empty",
			}
		}

		// Разделение заголовка на тип и токен
		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			return errorz.APIError{
				Status: http.StatusUnauthorized,
				Err:    fmt.Errorf("invalid authorization header"),
				Msg:    "invalid auth header",
			}
		}

		// Проверка токена и получение идентификатора пользователя
		claims, err := utils.VerifyToken(parts[1], []byte(h.jwtConfig.AccessTokenSecret()))
		if err != nil {
			return errorz.APIError{
				Status: http.StatusUnauthorized,
				Err:    err,
				Msg:    "invalid auth token",
			}
		}

		// Установка идентификатора пользователя и токена в контексте запроса

		md := metadata.New(map[string]string{"user_id": claims.UserID})
		ctx = metadata.NewOutgoingContext(r.Context(), md)
		ctx = context.WithValue(ctx, "user_id", claims.UserID)

		// Передача запроса следующему обработчику
		next.ServeHTTP(w, r.WithContext(ctx))
		return nil
	})
}
