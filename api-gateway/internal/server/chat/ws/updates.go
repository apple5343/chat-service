package ws

import (
	"encoding/json"
	"log"
	"slices"
)

type UserAdded struct {
	Type   string `json:"type",omitempty,default:"user_added"`
	UserID string `json:"user_id"`
	RoomID string `json:"room_id"`
}

type UserRemoved struct {
	Type   string `json:"type",omitempty,default:"user_removed"`
	UserID string `json:"user_id"`
	RoomID string `json:"room_id"`
}

func (h *Hub) ProcessUpdates() {
	for {
		select {
		case update, ok := <-h.Updates:
			if !ok {
				return
			}
			switch update.Type {
			case "user_added":
				var userAdded UserAdded
				err := json.Unmarshal([]byte(update.Data), &userAdded)
				if err != nil {
					log.Printf("error: %v", err)
					continue
				}
				h.AddUserToRoom(userAdded.RoomID, userAdded.UserID)
			case "user_removed":
				var userRemoved UserRemoved
				err := json.Unmarshal([]byte(update.Data), &userRemoved)
				if err != nil {
					log.Printf("error: %v", err)
					continue
				}
				h.RemoveUserFromRoom(userRemoved.RoomID, userRemoved.UserID)
			}
		case <-h.ctx.Done():
			return
		}
	}
}

func (h *Hub) RemoveUserFromRoom(roomID string, userID string) {
	if roomID == "" || userID == "" {
		return
	}
	client, ok := h.Users[userID]
	if !ok {
		return
	}
	room, ok := h.Rooms[roomID]
	if !ok {
		return
	}
	delete(room.Clients, userID)
	client.Updates <- UserAdded{
		Type:   "user_removed",
		UserID: userID,
		RoomID: roomID,
	}
	client.RoomIDs = slices.DeleteFunc(client.RoomIDs, func(id string) bool {
		return id == roomID
	})
}

func (h *Hub) AddUserToRoom(roomID string, userID string) {
	if roomID == "" || userID == "" {
		return
	}
	client, ok := h.Users[userID]
	if !ok {
		return
	}
	room, ok := h.Rooms[roomID]
	if !ok {
		room = &Room{
			ID:      roomID,
			Name:    roomID,
			Clients: map[string]*Client{},
		}
		h.Rooms[roomID] = room
	}
	room.Clients[userID] = client
	client.RoomIDs = append(client.RoomIDs, roomID)
	client.Updates <- UserAdded{
		Type:   "user_added",
		UserID: userID,
		RoomID: roomID,
	}
}
