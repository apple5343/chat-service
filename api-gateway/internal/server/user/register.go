package user

import (
	"api-gateway/internal/errorz"
	"api-gateway/internal/models/request"
	"api-gateway/internal/models/response"
	"api-gateway/internal/server/helper"
	userpb "api-gateway/pkg/api/user_v1"
	"encoding/json"
	"fmt"
	"net/http"
)

// Register обрабатывает запрос на регистрацию нового пользователя.
func (h *Handler) Register(w http.ResponseWriter, r *http.Request) error {
	// Получение контекста из запроса
	ctx := r.Context()

	// Декодируем тело запроса в структуру registerRequest
	var registerRequest request.Register
	if err := json.NewDecoder(r.Body).Decode(&registerRequest); err != nil {
		return errorz.APIError{
			Status: http.StatusBadRequest,
			Err:    fmt.Errorf("failed to decode Register request body: %w", err),
			Msg:    "error decoding request body",
		}
	}

	registerResp, err := h.userAPIClient.Register(ctx, &userpb.RegisterRequest{
		Email:    registerRequest.Email,
		Password: registerRequest.Password,
		Grade:    registerRequest.Grade,
	})
	if err != nil {
		return errorz.APIError{
			Status: http.StatusInternalServerError,
			Err:    fmt.Errorf("failed to register user: %w", err),
			Msg:    "failed to register user",
		}
	}

	resp := response.Register{ID: registerResp.GetId()}

	// Возвращаем успешный ответ со статусом 201 (Created)
	return helper.WriteJSON(w, http.StatusCreated, resp)
}
