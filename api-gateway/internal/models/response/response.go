package response

// Error ответ с ошибкой.
type Error struct {
	Error string `json:"error"`
}

// Message ответ с сообщением.
type Message struct {
	Message string `json:"message"`
}
