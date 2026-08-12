package entity

const (
	UserRemovedUpdate = "user_removed"
	UserAddedUpdate   = "user_added"
)

type Message struct {
	Type    string `json:"type,omitempty"`
	Content string `json:"content"`
	RoomID  string `json:"room_id"`
	UserID  string `json:"user_id"`
	Time    string `json:"time,omitempty"`
}

type Update struct {
	Type string `json:"type"`
	Data string `json:"data"`
}
