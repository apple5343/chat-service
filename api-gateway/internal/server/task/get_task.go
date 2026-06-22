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

func (h *Handler) GetTask(w http.ResponseWriter, r *http.Request) error {
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

	// Получаем id пользователя из url параметров
	taskID := chi.URLParam(r, "task_id")
	if len(taskID) == 0 {
		return errorz.APIError{
			Status: http.StatusBadRequest,
			Err:    fmt.Errorf("failed to get task id"),
			Msg:    "failed to get task id",
		}
	}

	getTaskByIDResp, err := h.taskAPIClient.GetTaskByID(ctx, &taskpb.GetTaskByIDRequest{
		Id: taskID,
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
				Msg:    "task not found",
			}
		case codes.Internal:
			return errorz.APIError{
				Status: http.StatusInternalServerError,
				Err:    err,
				Msg:    "failed to get task",
			}
		}
	}

	resp := response.GetTask{
		ID:          getTaskByIDResp.GetTask().GetId(),
		Title:       getTaskByIDResp.GetTask().GetTitle(),
		Description: getTaskByIDResp.GetTask().GetDescription(),
		Status:      getTaskByIDResp.GetTask().GetStatus(),
		ProjectID:   getTaskByIDResp.GetTask().GetProjectId(),
		ExecutorID:  getTaskByIDResp.GetTask().GetExecutorId(),
		Deadline:    getTaskByIDResp.GetTask().GetDeadline().AsTime(),
		CreatedAt:   getTaskByIDResp.GetTask().GetCreatedAt().AsTime(),
		UpdatedAt:   getTaskByIDResp.GetTask().GetUpdatedAt().AsTime(),
	}

	return helper.WriteJSON(w, http.StatusOK, resp)
}
