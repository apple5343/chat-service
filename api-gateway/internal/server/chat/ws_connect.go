package chat

import (
	"api-gateway/internal/errorz"
	ws "api-gateway/internal/server/chat/ws"
	chatpb "api-gateway/pkg/api/chat_v1"
	"fmt"
	"net/http"

	"github.com/gorilla/websocket"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin:     func(r *http.Request) bool { return true }, //TODO нужно будет убрать
}

func (h *Handler) JoinRooms(w http.ResponseWriter, r *http.Request) error {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		return errorz.APIError{
			Status: http.StatusBadRequest,
			Err:    err,
			Msg:    "failed to upgrade connection",
		}
	}
	userID, ok := r.Context().Value("user_id").(string)
	fmt.Println("Conected: ", userID)
	if !ok {
		return errorz.APIError{
			Status: http.StatusBadRequest,
			Err:    fmt.Errorf("failed to get user id"),
			Msg:    "failed to get user id",
		}
	}
	userChats, err := h.chatAPIClient.GetUserChats(r.Context(), &chatpb.GetUserChatsRequest{
		UserId: userID,
	})
	fmt.Println("User Chats: ", userChats.GetProjectIds())
	if err != nil {
		fmt.Println(err)
		st, _ := status.FromError(err)
		switch st.Code() {
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
				Msg:    "failed to get user chats",
			}
		}
	}
	client := &ws.Client{
		ID:      userID,
		Conn:    conn,
		Updates: make(chan interface{}),
		RoomIDs: userChats.GetProjectIds(),
	}
	fmt.Println(client)
	h.hub.Register <- client
	fmt.Println("Client register")
	go client.ReadMessage(h.hub)
	go client.WriteUpdates()
	return nil
}
