package client

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/gorilla/websocket"
)

type Message struct {
	RoomID  string `json:"room_id"`
	Content string `json:"content"`
}

func (c *Client) ConnectWS(ctx context.Context, token string) (*websocket.Conn, error) {
	header := http.Header{}
	header.Add("Authorization", "Bearer "+token)
	conn, _, err := websocket.DefaultDialer.DialContext(ctx, c.wsUrl+"/chats/ws", header)

	return conn, err
}

func (c *Client) SendMessage(conn *websocket.Conn, message *Message) error {
	body, err := json.Marshal(message)
	if err != nil {
		return err
	}
	return conn.WriteMessage(websocket.TextMessage, body)
}
