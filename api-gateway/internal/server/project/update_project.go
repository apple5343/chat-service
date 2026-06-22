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

	"github.com/go-chi/chi/v5"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// UpdateProject обрабатывает запрос на обновление проекта.
func (h *Handler) UpdateProject(w http.ResponseWriter, r *http.Request) error {
	// Получение контекста
	ctx := r.Context()

	// Получаем id пользователя из url параметров
	projectID := chi.URLParam(r, "project_id")
	if len(projectID) == 0 {
		return errorz.APIError{
			Status: http.StatusBadRequest,
			Err:    fmt.Errorf("failed to get project id"),
			Msg:    "failed to get project id",
		}
	}

	// Декодируем тело запроса в структуру getAccessTokenRequest
	var updateProjectRequest request.UpdateProject
	if err := json.NewDecoder(r.Body).Decode(&updateProjectRequest); err != nil {
		return errorz.APIError{
			Status: http.StatusBadRequest,
			Err:    err,
			Msg:    "failed to decode request body",
		}
	}

	// Получение ответа от api клиента
	_, err := h.projectAPIClient.UpdateProject(ctx, &projectpb.UpdateProjectRequest{
		Id:          projectID,
		Title:       updateProjectRequest.Title,
		Description: updateProjectRequest.Description,
		AdminId:     updateProjectRequest.AdminID,
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
		case codes.PermissionDenied:
			return errorz.APIError{
				Status: http.StatusForbidden,
				Err:    err,
				Msg:    "permission denied",
			}
		case codes.NotFound:
			return errorz.APIError{
				Status: http.StatusNotFound,
				Err:    err,
				Msg:    "project not found",
			}
		case codes.Internal:
			return errorz.APIError{
				Status: http.StatusInternalServerError,
				Err:    err,
				Msg:    "failed to update project",
			}
		}
	}

	// Обновляем чат
	_, err = h.chatApiClient.UpdateChat(ctx, &chatpb.UpdateChatRequest{
		Chat: &chatpb.Chat{
			ChatId: projectID,
			Name:   updateProjectRequest.Title,
		},
	})

	// Возвращаем ответ
	return helper.WriteJSON(w, http.StatusOK, response.Message{Message: "project updated successfully"})
}
