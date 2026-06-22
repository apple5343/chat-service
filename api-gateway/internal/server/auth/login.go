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

// Login обрабатывает запрос на вход пользователя.
func (h *Handler) Login(w http.ResponseWriter, r *http.Request) error {
	// Получение контекста из запроса
	ctx := r.Context()

	// Декодируем тело запроса в структуру loginRequest
	var loginRequest request.Login
	if err := json.NewDecoder(r.Body).Decode(&loginRequest); err != nil {
		return errorz.APIError{
			Status: http.StatusBadRequest,
			Err:    fmt.Errorf("failed to decode Login request body: %w", err),
			Msg:    "error decoding request body",
		}
	}

	// Получение ответа от api клиента
	loginResp, err := h.authAPIClient.Login(ctx, &authpb.LoginRequest{
		Email:    loginRequest.Email,
		Password: loginRequest.Password,
	})
	if err != nil {
		return errorz.APIError{
			Status: http.StatusUnauthorized,
			Err:    fmt.Errorf("login failed: %w", err),
			Msg:    "login failed",
		}
	}

	// Формируем и возвращаем ответ
	resp := response.Login{
		RefreshToken: loginResp.GetRefreshToken(),
	}
	return helper.WriteJSON(w, http.StatusOK, resp)
}
