package ws

import (
	"api-gateway/internal/entity"
	"encoding/json"
	"errors"
	"log"
	"sync"
	"time"

	"github.com/gorilla/websocket"
)

type User struct {
	id      string
	conn    *websocket.Conn
	rooms   []string
	roomsMu sync.RWMutex
	sendCh  chan []byte
	done    chan struct{}
	once    sync.Once
}

func NewUser(id string, conn *websocket.Conn, rooms []string) *User {
	return &User{
		id:     id,
		conn:   conn,
		rooms:  rooms,
		sendCh: make(chan []byte, 128),
		done:   make(chan struct{}),
		once:   sync.Once{},
	}
}

func (u *User) AddRoom(roomID string) {
	u.roomsMu.Lock()
	defer u.roomsMu.Unlock()

	for _, id := range u.rooms {
		if id == roomID {
			return
		}
	}

	u.rooms = append(u.rooms, roomID)
}

func (u *User) RemoveRoom(roomID string) {
	u.roomsMu.Lock()
	defer u.roomsMu.Unlock()

	rooms := make([]string, 0, len(u.rooms))
	for _, id := range u.rooms {
		if id != roomID {
			rooms = append(rooms, id)
		}
	}

	u.rooms = rooms
}

func (u *User) Rooms() []string {
	u.roomsMu.RLock()
	defer u.roomsMu.RUnlock()
	return append([]string(nil), u.rooms...)
}

func (u *User) SendMessage(message []byte) error {
	select {
	case u.sendCh <- message:
		return nil
	case <-u.done:
		return errors.New("user closed")
	default:
		return errors.New("user send buffer full")
	}
}

func (u *User) CloseConn() {
	u.once.Do(func() {
		close(u.done)
		u.conn.Close()
	})
}

func (u *User) StartRead(sendCh chan<- *entity.Message) {
	u.conn.SetReadDeadline(time.Now().Add(60 * time.Second))
	u.conn.SetPongHandler(func(appData string) error {
		u.conn.SetReadDeadline(time.Now().Add(60 * time.Second))
		return nil
	})

	for {
		_, message, err := u.conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Println("ws read error: ", err)
			}
			break
		}
		var userMessage entity.Message
		err = json.Unmarshal(message, &userMessage)
		if err != nil {
			continue
		}
		userMessage.UserID = u.id
		select {
		case sendCh <- &userMessage:
		case <-u.done:
			return
		default:
			log.Println("error: too many requests")
		}
	}
}

func (u *User) WritePump() {
	ticker := time.NewTicker(30 * time.Second)
	defer ticker.Stop()
	defer u.CloseConn()

	for {
		select {
		case data, ok := <-u.sendCh:
			if !ok {
				return
			}

			_ = u.conn.SetWriteDeadline(time.Now().Add(10 * time.Second))
			if err := u.conn.WriteMessage(websocket.TextMessage, data); err != nil {
				return
			}

		case <-ticker.C:
			err := u.conn.WriteControl(
				websocket.PingMessage,
				nil,
				time.Now().Add(time.Second),
			)
			if err != nil {
				return
			}

		case <-u.done:
			return
		}
	}
}
