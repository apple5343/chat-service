package user

import (
	"api-gateway/internal/errorz"
	"api-gateway/internal/models/response"
	"api-gateway/internal/server/helper"
	userpb "api-gateway/pkg/api/user_v1"
	"fmt"
	"net/http"

	"github.com/go-chi/chi/v5"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (h *Handler) DeleteUser(w http.ResponseWriter, r *http.Request) error {
	ctx := r.Context()

	// Получаем id пользователя из url параметров
	userID := chi.URLParam(r, "user_id")
	if len(userID) == 0 {
		return errorz.APIError{
			Status: http.StatusBadRequest,
			Err:    fmt.Errorf("failed to get user id"),
			Msg:    "failed to get user id",
		}
	}

	_, err := h.userAPIClient.DeleteUser(ctx, &userpb.DeleteUserRequest{
		Id: userID,
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
				Msg:    "user not found",
			}
		case codes.Internal:
			return errorz.APIError{
				Status: http.StatusInternalServerError,
				Err:    err,
				Msg:    "failed to delete user",
			}
		}
	}

	return helper.WriteJSON(w, http.StatusNoContent, response.Message{Message: "user deleted successfully"})
}
