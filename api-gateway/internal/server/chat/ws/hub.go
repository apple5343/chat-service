package ws

import (
	"api-gateway/internal/entity"
	"context"
	"fmt"
)

type ChatService interface {
	WriteMessage(msg *entity.Message)
	ReadMessage(chan *entity.Message)
	ReadUpdates(chan *entity.Update)
}

type Room struct {
	ID      string `json:"id"`
	Name    string `json:"name"`
	Clients map[string]*Client
}

type Hub struct {
	Rooms      map[string]*Room
	Register   chan *Client
	Unregister chan *Client
	Broadcast  chan *entity.Message
	Send       chan *entity.Message
	Updates    chan *entity.Update
	Users      map[string]*Client
	Chat       ChatService
	ctx        context.Context
}

func NewHub(chat ChatService) *Hub {

	broadCastch := make(chan *entity.Message, 10)
	chat.ReadMessage(broadCastch)

	updates := make(chan *entity.Update, 10)
	chat.ReadUpdates(updates)

	return &Hub{
		Rooms:      map[string]*Room{},
		Register:   make(chan *Client),
		Unregister: make(chan *Client),
		Users:      map[string]*Client{},
		Send:       make(chan *entity.Message),
		Chat:       chat,
		Broadcast:  broadCastch,
		Updates:    updates,
		ctx:        context.Background(),
	}
}

func (h *Hub) Worker(ctx context.Context) {
	for {
		select {
		case cl := <-h.Register:
			for _, room := range cl.RoomIDs {
				r, ok := h.Rooms[room]
				if !ok {
					r = &Room{
						ID:      room,
						Name:    room,
						Clients: map[string]*Client{},
					}
					h.Rooms[room] = r
				}
				r.Clients[cl.ID] = cl
			}
			h.Users[cl.ID] = cl

		case cl := <-h.Unregister:
			for _, room := range cl.RoomIDs {
				r, ok := h.Rooms[room]
				if !ok {
					continue
				}
				delete(r.Clients, cl.ID)
			}
			delete(h.Users, cl.ID)

		case msg := <-h.Broadcast:
			fmt.Println(msg)
			room, ok := h.Rooms[msg.RoomID]
			if !ok {
				continue
			}
			for _, cl := range room.Clients {
				cl.Updates <- msg
			}

		case msg := <-h.Send:
			h.Chat.WriteMessage(msg)

		case <-ctx.Done():
			return
		}
	}
}

func (h *Hub) Run() {
	for i := 0; i < 10; i++ {
		go h.Worker(h.ctx)
		go h.ProcessUpdates()
	}
}

func (h *Hub) Shutdown() {
	h.ctx.Done()
	close(h.Register)
	close(h.Unregister)
	close(h.Broadcast)
	close(h.Send)
	for _, room := range h.Rooms {
		for _, cl := range room.Clients {
			close(cl.Updates)
		}
	}
}
