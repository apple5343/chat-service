package project

import (
	"api-gateway/internal/errorz"
	"api-gateway/internal/models/response"
	"api-gateway/internal/server/helper"
	"net/http"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// ListUserProjects обрабатывает запрос на получения списка проекта пользователя.
func (h *Handler) ListUserProjects(w http.ResponseWriter, r *http.Request) error {
	// Получение контекста
	ctx := r.Context()

	// Получение ответа от api клиента
	listUserProjectsResp, err := h.projectAPIClient.ListUserProjects(ctx, nil)
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
				Msg:    "failed to list user projects",
			}
		}
	}

	// Формируем и возвращаем ответ
	resp := make([]response.GetProject, len(listUserProjectsResp.GetProject()))
	for i, project := range listUserProjectsResp.GetProject() {
		resp[i] = response.GetProject{
			ID:                           project.GetId(),
			Title:                        project.GetTitle(),
			Description:                  project.GetDescription(),
			Users:                        project.GetUsers(),
			AdminID:                      project.GetAdminId(),
			NotificationSubscribersTGIds: project.GetNotificationSubscribersTgIds(),
			CreatedAt:                    project.GetCreatedAt().AsTime(),
		}
	}
	return helper.WriteJSON(w, http.StatusOK, response.ListUserProjects{Projects: resp})
}
