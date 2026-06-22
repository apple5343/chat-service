package project

import (
	"api-gateway/internal/errorz"
	"api-gateway/internal/models/request"
	"api-gateway/internal/models/response"
	"api-gateway/internal/server/helper"
	chatpb "api-gateway/pkg/api/chat_v1"
	projectpb "api-gateway/pkg/api/project_v1"
	"encoding/json"
	"fmt"
	"net/http"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// CreateProject обрабатывает запрос на создание проекта.
func (h *Handler) CreateProject(w http.ResponseWriter, r *http.Request) error {
	// Получение контекста
	ctx := r.Context()

	userID, ok := r.Context().Value("user_id").(string)
	if !ok {
		return errorz.APIError{
			Status: http.StatusUnauthorized,
			Err:    fmt.Errorf("failed to get user id"),
			Msg:    "failed to get user id",
		}
	}

	// Декодируем тело запроса в структуру createProjectRequest
	var createProjectRequest request.CreateProject
	if err := json.NewDecoder(r.Body).Decode(&createProjectRequest); err != nil {
		return errorz.APIError{
			Status: http.StatusBadRequest,
			Err:    fmt.Errorf("failed to decode CreateProject request body: %w", err),
			Msg:    "error decoding request body",
		}
	}
	createProjectRequest.Users = append(createProjectRequest.Users, userID)

	// Получение ответа от api клиента
	createProjectResp, err := h.projectAPIClient.CreateProject(ctx, &projectpb.CreateProjectRequest{
		Title:       createProjectRequest.Title,
		Description: createProjectRequest.Description,
		Users:       createProjectRequest.Users,
	})
	if err != nil {
		st, _ := status.FromError(err)

		switch st.Code() {
		case codes.Unauthenticated:
			return errorz.APIError{
				Status: http.StatusUnauthorized,
				Err:    err,
				Msg:    "authentication required",
			}
		case codes.Internal:
			return errorz.APIError{
				Status: http.StatusInternalServerError,
				Err:    err,
				Msg:    "failed to create project",
			}
		}
	}

	// Создание нового чата
	_, err = h.chatApiClient.CreateChat(ctx, &chatpb.CreateChatRequest{
		ProjectId: createProjectResp.GetId(),
		Name:      createProjectRequest.Title,
		Member:    createProjectRequest.Users,
	})

	// Формируем и возвращаем ответ
	resp := response.CreateProject{
		ID: createProjectResp.GetId(),
	}
	return helper.WriteJSON(w, http.StatusCreated, resp)
}
