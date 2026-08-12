package storage

import (
	"sync"
	"time"
)

type MessagePath struct {
	SendedAt   time.Time
	ReceivedAt time.Time
	Received   bool
}

type Storage struct {
	m  map[string]*MessagePath
	mu sync.Mutex
}

func NewStorage() *Storage {
	return &Storage{
		m:  make(map[string]*MessagePath),
		mu: sync.Mutex{},
	}
}

func (s *Storage) Send(id string, at time.Time) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.m[id] = &MessagePath{SendedAt: at}
}

func (s *Storage) Receive(id string, at time.Time) {
	s.mu.Lock()
	defer s.mu.Unlock()
	message, ok := s.m[id]
	if !ok {
		return
	}
	message.ReceivedAt = at
	message.Received = true
}

func (s *Storage) Messages() []MessagePath {
	s.mu.Lock()
	defer s.mu.Unlock()

	result := make([]MessagePath, 0, len(s.m))
	for _, msg := range s.m {
		result = append(result, MessagePath{SendedAt: msg.SendedAt, ReceivedAt: msg.ReceivedAt, Received: msg.Received})
	}
	return result
}
