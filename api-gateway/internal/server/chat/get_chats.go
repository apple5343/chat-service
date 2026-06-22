package chat

import (
	"api-gateway/internal/errorz"
	"api-gateway/internal/server/helper"
	chatpb "api-gateway/pkg/api/chat_v1"
	"net/http"
)

func (h *Handler) GetChats(w http.ResponseWriter, r *http.Request) error {
	userID, _ := r.Context().Value("user_id").(string)

	getChatsResp, err := h.chatAPIClient.GetUserChats(r.Context(), &chatpb.GetUserChatsRequest{
		UserId: userID,
	})
	if err != nil {
		return errorz.APIError{
			Status: http.StatusInternalServerError,
			Err:    err,
			Msg:    "failed to get chats",
		}
	}

	return helper.WriteJSON(w, http.StatusOK, getChatsResp)
}
