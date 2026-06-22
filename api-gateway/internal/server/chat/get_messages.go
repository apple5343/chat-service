package chat

import (
	"api-gateway/internal/entity"
	"api-gateway/internal/errorz"
	"api-gateway/internal/server/helper"
	chatpb "api-gateway/pkg/api/chat_v1"
	"fmt"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (h *Handler) GetMessages(w http.ResponseWriter, r *http.Request) error {
	projectID := chi.URLParam(r, "project_id")
	if len(projectID) == 0 {
		return errorz.APIError{
			Status: http.StatusBadRequest,
			Err:    fmt.Errorf("failed to get project id"),
			Msg:    "failed to get project id",
		}
	}
	cursorStr := r.URL.Query().Get("cursor")
	cursor, err := strconv.Atoi(cursorStr)
	if err != nil {
		return errorz.APIError{
			Status: http.StatusBadRequest,
			Err:    err,
			Msg:    "failed to parse cursor",
		}
	}
	limitStr := r.URL.Query().Get("limit")
	limit, err := strconv.Atoi(limitStr)
	if err != nil {
		return errorz.APIError{
			Status: http.StatusBadRequest,
			Err:    err,
			Msg:    "failed to parse limit",
		}
	}
	getMessagesResp, err := h.chatAPIClient.GetMessages(r.Context(), &chatpb.GetMessagesRequest{
		ProjectId: projectID,
		Cursor:    int32(cursor),
		Limit:     int32(limit),
		UserId:    r.Context().Value("user_id").(string),
	})
	if err != nil {
		st, _ := status.FromError(err)
		switch st.Code() {
		case codes.NotFound:
			return errorz.APIError{
				Status: http.StatusNotFound,
				Err:    err,
				Msg:    "chat not found",
			}
		case codes.PermissionDenied:
			return errorz.APIError{
				Status: http.StatusForbidden,
				Err:    err,
				Msg:    "permission denied",
			}
		case codes.InvalidArgument:
			return errorz.APIError{
				Status: http.StatusBadRequest,
				Err:    err,
				Msg:    "invalid argument",
			}
		case codes.Internal:
			return errorz.APIError{
				Status: http.StatusInternalServerError,
				Err:    err,
				Msg:    "failed to get messages",
			}
		}
	}

	resp := []*entity.Message{}
	for _, message := range getMessagesResp.Messages {
		resp = append(resp, &entity.Message{
			RoomID:  message.GetProjectId(),
			UserID:  message.GetUserId(),
			Content: message.GetContent(),
			Time:    message.GetTimestamp(),
		})
	}

	return helper.WriteJSON(w, http.StatusOK, resp)
}
