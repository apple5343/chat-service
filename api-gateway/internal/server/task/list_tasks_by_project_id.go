package task

import (
	"api-gateway/internal/errorz"
	"api-gateway/internal/models/response"
	"api-gateway/internal/server/helper"
	taskpb "api-gateway/pkg/api/task_v1"
	"fmt"
	"net/http"

	"github.com/go-chi/chi/v5"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (h *Handler) ListTasksByProjectID(w http.ResponseWriter, r *http.Request) error {
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

	listTasksByProjectIDResp, err := h.taskAPIClient.ListTasksByProjectID(ctx, &taskpb.ListTasksByProjectIDRequest{
		ProjectId: projectID,
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
				Msg:    "failed to list tasks by project id",
			}
		}
	}

	resp := make([]response.GetTask, len(listTasksByProjectIDResp.Tasks))
	for i, task := range listTasksByProjectIDResp.Tasks {
		resp[i] = response.GetTask{
			ID:          task.GetId(),
			Title:       task.GetTitle(),
			Description: task.GetDescription(),
			Status:      task.GetStatus(),
			ProjectID:   task.GetProjectId(),
			ExecutorID:  task.GetExecutorId(),
			Deadline:    task.GetDeadline().AsTime(),
			CreatedAt:   task.GetCreatedAt().AsTime(),
			UpdatedAt:   task.GetUpdatedAt().AsTime(),
		}
	}

	return helper.WriteJSON(w, http.StatusOK, response.ListTasksByProjectID{Tasks: resp})
}
