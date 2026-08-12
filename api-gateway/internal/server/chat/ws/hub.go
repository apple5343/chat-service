package ws

import (
	"api-gateway/internal/entity"
	"context"
	"log"
	"sync"
)

type ChatService interface {
	WriteMessage(context.Context, *entity.Message) error
	ReadMessages(context.Context, chan *entity.Message)
	ReadUpdates(context.Context, chan *entity.Update)
}

type Hub struct {
	rooms         map[string]*Room
	roomsMu       sync.RWMutex
	users         map[string]*User
	usersMu       sync.Mutex
	messagesInCh  chan *entity.Message
	messagesOutCh chan *entity.Message
	updatesCh     chan *entity.Update
	connectionsWg sync.WaitGroup
	wg            sync.WaitGroup
	chat          ChatService
	ctx           context.Context
	cancel        func()
}

func NewHub(chat ChatService) *Hub {
	return &Hub{
		rooms:         make(map[string]*Room),
		roomsMu:       sync.RWMutex{},
		users:         make(map[string]*User),
		usersMu:       sync.Mutex{},
		messagesInCh:  make(chan *entity.Message, 2000),
		messagesOutCh: make(chan *entity.Message, 2000),
		updatesCh:     make(chan *entity.Update, 2000),
		connectionsWg: sync.WaitGroup{},
		wg:            sync.WaitGroup{},
		chat:          chat,
	}
}

func (h *Hub) Run() {
	h.ctx, h.cancel = context.WithCancel(context.Background())

	h.wg.Add(5)
	go func() {
		defer h.wg.Done()
		h.chat.ReadMessages(h.ctx, h.messagesInCh)
	}()

	go func() {
		defer h.wg.Done()
		h.chat.ReadUpdates(h.ctx, h.updatesCh)
	}()

	go func() {
		defer h.wg.Done()
		h.ProcessUpdates(h.ctx)
	}()

	go func() {
		defer h.wg.Done()
		h.ProcessMessagesIn(h.ctx)
	}()

	go func() {
		defer h.wg.Done()
		h.ProcessMessagesOut(h.ctx)
	}()
}

func (h *Hub) ConnectUser(user *User) {
	h.usersMu.Lock()
	if old, ok := h.users[user.id]; ok && old != user {
		old.CloseConn()
	}
	h.users[user.id] = user
	h.usersMu.Unlock()

	h.roomsMu.RLock()
	unavailableRooms := []string{}
	for _, roomID := range user.Rooms() {
		if room, ok := h.rooms[roomID]; ok {
			room.AddUser(user)
		} else {
			unavailableRooms = append(unavailableRooms, roomID)
		}
	}
	h.roomsMu.RUnlock()

	h.connectionsWg.Add(2)

	go func() {
		defer h.connectionsWg.Done()
		user.WritePump()
	}()

	go func() {
		defer h.connectionsWg.Done()
		user.StartRead(h.messagesOutCh)
		h.DisconnectUser(user)
	}()
	if len(unavailableRooms) == 0 {
		return
	}

	h.roomsMu.Lock()
	for _, roomID := range unavailableRooms {
		room, ok := h.rooms[roomID]
		if !ok {
			room = NewRoom(roomID)
			h.rooms[roomID] = room
		}
		room.AddUser(user)
	}
	h.roomsMu.Unlock()
}

func (h *Hub) removeRoom(roomID string) {
	h.roomsMu.Lock()
	defer h.roomsMu.Unlock()

	room, ok := h.rooms[roomID]
	if ok && room.IsEmpty() {
		delete(h.rooms, roomID)
	}
}

func (h *Hub) DisconnectUser(user *User) {
	h.usersMu.Lock()
	current, ok := h.users[user.id]
	if !ok || current != user {
		h.usersMu.Unlock()
		return
	}
	delete(h.users, user.id)
	h.usersMu.Unlock()
	user.CloseConn()

	h.roomsMu.RLock()
	emptyRooms := []string{}
	for _, roomID := range user.Rooms() {
		room, ok := h.rooms[roomID]
		if !ok {
			log.Println("room missing")
			continue
		}
		room.RemoveUser(user)
		if room.IsEmpty() {
			emptyRooms = append(emptyRooms, roomID)
		}
	}
	h.roomsMu.RUnlock()
	if len(emptyRooms) == 0 {
		return
	}

	h.roomsMu.Lock()
	defer h.roomsMu.Unlock()
	for _, roomID := range emptyRooms {
		room, ok := h.rooms[roomID]
		if ok && room.IsEmpty() {
			delete(h.rooms, roomID)
		}
	}
}

func (h *Hub) CleanConnections() {
	h.usersMu.Lock()
	defer h.usersMu.Unlock() //Lock не дает подключать новый пользователей
	for _, user := range h.users {
		user.CloseConn()
	}
}

func (h *Hub) Shutdown() {
	h.CleanConnections()
	h.connectionsWg.Wait()

	if h.cancel != nil {
		h.cancel()
	}

	h.wg.Wait()
}
