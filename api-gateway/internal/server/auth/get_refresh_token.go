package auth

import (
	"api-gateway/internal/errorz"
	"api-gateway/internal/models/request"
	"api-gateway/internal/models/response"
	"api-gateway/internal/server/helper"
	authpb "api-gateway/pkg/api/auth_v1"
	"encoding/json"
	"fmt"
	"net/http"
)

// GetRefreshToken обрабатывает запрос на получения refresh токена.
func (h *Handler) GetRefreshToken(w http.ResponseWriter, r *http.Request) error {
	// Получение контекста
	ctx := r.Context()

	// Декодируем тело запроса в структуру getRefreshTokenRequest
	var getRefreshTokenRequest request.GetRefreshToken
	if err := json.NewDecoder(r.Body).Decode(&getRefreshTokenRequest); err != nil {
		return errorz.APIError{
			Status: http.StatusBadRequest,
			Err:    fmt.Errorf("failed to decode GetRefreshToken request body: %w", err),
			Msg:    "error decoding request body",
		}
	}

	// Получение ответа от api клиента
	getRefreshTokenResp, err := h.authAPIClient.GetRefreshToken(ctx, &authpb.GetRefreshTokenRequest{
		RefreshToken: getRefreshTokenRequest.RefreshToken,
	})
	if err != nil {
		return errorz.APIError{
			Status: http.StatusUnauthorized,
			Err:    fmt.Errorf("failed to get refresh token: %w", err),
			Msg:    "failed to get refresh token",
		}
	}

	// Формируем и возвращаем ответ
	resp := response.GetRefreshToken{
		RefreshToken: getRefreshTokenResp.GetRefreshToken(),
	}
	return helper.WriteJSON(w, http.StatusOK, resp)
}
