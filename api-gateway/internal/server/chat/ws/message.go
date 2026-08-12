package ws

import (
	"api-gateway/internal/entity"
	"context"
	"encoding/json"
	"log"
	"sync"

	"golang.org/x/sync/semaphore"
)

func (h *Hub) ProcessMessagesOut(ctx context.Context) {
	sem := semaphore.NewWeighted(1000)
	wg := sync.WaitGroup{}
	defer func() {
		wg.Wait()
	}()

	for {
		select {
		case message, ok := <-h.messagesOutCh:
			if !ok {
				return
			}
			if message == nil {
				continue
			}
			if !h.canSendMessage(message) {
				continue
			}

			if err := sem.Acquire(ctx, 1); err != nil {
				return
			}

			wg.Add(1)
			go func() {
				defer wg.Done()
				defer sem.Release(1)

				if err := h.chat.WriteMessage(ctx, message); err != nil {
					log.Println("error write message: ", err.Error())
				}
			}()
		case <-ctx.Done():
			return
		}
	}
}

func (h *Hub) ProcessMessagesIn(ctx context.Context) {
	sem := semaphore.NewWeighted(1000)
	wg := sync.WaitGroup{}
	defer func() {
		wg.Wait()
	}()

	for {
		select {
		case message, ok := <-h.messagesInCh:
			if !ok {
				return
			}
			if message == nil {
				continue
			}
			if err := sem.Acquire(ctx, 1); err != nil {
				return
			}

			wg.Add(1)
			go func() {
				defer wg.Done()
				defer sem.Release(1)

				if err := h.ProcessMessage(message); err != nil {
					log.Println("error procees incoming message: ", err.Error())
				}
			}()
		case <-ctx.Done():
			return
		}
	}
}

func (h *Hub) ProcessMessage(message *entity.Message) error {
	if message == nil {
		return nil
	}
	h.roomsMu.RLock()
	room, ok := h.rooms[message.RoomID]
	h.roomsMu.RUnlock()
	if !ok {
		return nil
	}
	content, err := json.Marshal(message)
	if err != nil {
		return err
	}
	room.SendMessage(content)
	return nil
}

func (h *Hub) canSendMessage(message *entity.Message) bool {
	if message == nil {
		return false
	}

	if message.UserID == "" || message.RoomID == "" {
		return false
	}

	h.usersMu.Lock()
	user, ok := h.users[message.UserID]
	h.usersMu.Unlock()

	if !ok {
		return false
	}

	h.roomsMu.RLock()
	room, ok := h.rooms[message.RoomID]
	h.roomsMu.RUnlock()

	if !ok {
		return false
	}

	return room.HasUser(user)
}
