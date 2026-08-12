package ws

import (
	"api-gateway/internal/entity"
	"context"
	"encoding/json"
	"fmt"
	"log"
)

type UserAddedUpdate struct {
	UserID string `json:"user_id"`
	RoomID string `json:"room_id"`
}

type UserRemovedUpdate struct {
	UserID string `json:"user_id"`
	RoomID string `json:"room_id"`
}

type UpdateMessage struct {
	Type    string                 `json:"type"`
	Details map[string]interface{} `json:"details"`
}

func (h *Hub) ProcessUpdates(ctx context.Context) {
	for {
		select {
		case update, ok := <-h.updatesCh:
			if !ok {
				return
			}
			if err := h.HandleUpdate(update); err != nil {
				log.Println("error handle update: ", err.Error())
			}
		case <-ctx.Done():
			return
		}
	}
}

func (h *Hub) HandleUpdate(update *entity.Update) error {
	switch update.Type {
	case entity.UserRemovedUpdate:
		return h.HandleUserRemoved(update)
	case entity.UserAddedUpdate:
		return h.HandleUserAdded(update)
	}
	return fmt.Errorf("invalid update type")
}

func (h *Hub) HandleUserRemoved(update *entity.Update) error {
	var data UserRemovedUpdate
	if err := json.Unmarshal([]byte(update.Data), &data); err != nil {
		return err
	}

	h.usersMu.Lock()
	user, ok := h.users[data.UserID]
	h.usersMu.Unlock()
	if !ok {
		return nil
	}

	h.roomsMu.RLock()
	room, ok := h.rooms[data.RoomID]
	h.roomsMu.RUnlock()
	if !ok {
		return nil
	}

	if !room.HasUser(user) {
		return nil
	}

	user.RemoveRoom(room.id)
	room.RemoveUser(user)

	if room.IsEmpty() {
		h.removeRoom(room.id)
	}

	message, err := json.Marshal(UpdateMessage{
		Type:    "removed_from_room",
		Details: map[string]interface{}{"room_id": data.RoomID},
	})
	if err != nil {
		return err
	}
	return user.SendMessage(message)
}

func (h *Hub) HandleUserAdded(update *entity.Update) error {
	var data UserAddedUpdate
	if err := json.Unmarshal([]byte(update.Data), &data); err != nil {
		return err
	}

	h.usersMu.Lock()
	user, ok := h.users[data.UserID]
	h.usersMu.Unlock()
	if !ok {
		return nil
	}

	h.roomsMu.Lock()
	room, ok := h.rooms[data.RoomID]
	if !ok {
		room = NewRoom(data.RoomID)
		h.rooms[data.RoomID] = room
	}
	h.roomsMu.Unlock()

	user.AddRoom(data.RoomID)
	room.AddUser(user)

	message, err := json.Marshal(UpdateMessage{
		Type:    "added_to_room",
		Details: map[string]interface{}{"room_id": data.RoomID},
	})
	if err != nil {
		return err
	}
	return user.SendMessage(message)
}
