package ws

import (
	"api-gateway/internal/entity"
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/gorilla/websocket"
)

type Client struct {
	ID      string `json:"id"`
	Conn    *websocket.Conn
	Updates chan interface{}
	RoomIDs []string `json:"roomIds"`
}

func (c *Client) WriteUpdates() {
	defer c.Conn.Close()

	for {
		update, ok := <-c.Updates
		if !ok {
			c.Conn.WriteMessage(websocket.CloseMessage, []byte{})
			c.Conn.Close()
			return
		}
		fmt.Println(update)
		c.Conn.WriteJSON(update)
	}
}

func (c *Client) ReadMessage(h *Hub) {
	fmt.Println("Start read message ", c.ID)
	defer func() {
		h.Unregister <- c
		c.Conn.Close()
	}()

	c.Conn.SetReadDeadline(time.Now().Add(time.Second * 60))
	c.Conn.SetPongHandler(func(string) error {
		c.Conn.SetReadDeadline(time.Now().Add(time.Second * 60))
		return nil
	})

	ticker := time.NewTicker(time.Second * 30)
	defer ticker.Stop()

	go func() {
		for range ticker.C {
			c.Conn.WriteMessage(websocket.PingMessage, []byte{})
		}
	}()

	for {
		_, m, err := c.Conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Printf("error: %v", err)
			}
			break
		}
		fmt.Println("Message: ", string(m), " ID: ", c.ID)
		var msg *entity.Message
		err = json.Unmarshal(m, &msg)
		if err != nil {
			log.Printf("error: %v", err)
			continue
		}
		msg.UserID = c.ID
		h.Send <- msg
	}
}
