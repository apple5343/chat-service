package ws

import (
	"log"
	"sync"
)

type Room struct {
	id    string
	users map[string]*User
	mu    sync.RWMutex
}

func NewRoom(roomID string) *Room {
	return &Room{
		id:    roomID,
		users: make(map[string]*User),
		mu:    sync.RWMutex{},
	}
}

func (r *Room) AddUser(user *User) {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.users[user.id] = user
}

func (r *Room) RemoveUser(user *User) {
	r.mu.Lock()
	defer r.mu.Unlock()
	if r.users[user.id] == user {
		delete(r.users, user.id)
	}
}

func (r *Room) SendMessage(message []byte) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	for _, user := range r.users {
		if err := user.SendMessage(message); err != nil {
			log.Println("send message: ", err.Error())
		}
	}
}

func (r *Room) HasUser(user *User) bool {
	r.mu.RLock()
	defer r.mu.RUnlock()
	return r.users[user.id] == user
}

func (r *Room) IsEmpty() bool {
	r.mu.RLock()
	defer r.mu.RUnlock()
	return len(r.users) == 0
}
