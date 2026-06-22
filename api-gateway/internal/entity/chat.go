package entity

type Message struct {
	Type    string `json:"type",omitempty,default:"message"`
	Content string `json:"content"`
	RoomID  string `json:"room_id"`
	UserID  string `json:"user_id"`
	Time    string `json:"time,omitempty"`
}

type Update struct {
	Type string `json:"type"`
	Data string `json:"data"`
}
