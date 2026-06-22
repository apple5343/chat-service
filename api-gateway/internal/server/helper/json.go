package helper

import (
	"encoding/json"
	"net/http"
)

// WriteJSON записывает данные в формате JSON в ответ на HTTP-запрос с указанным статус-кодом.
// Она устанавливает статус ответа, кодирует данные в JSON и отправляет их в тело ответа.
func WriteJSON(w http.ResponseWriter, code int, data any) error {
	w.WriteHeader(code)
	return json.NewEncoder(w).Encode(data)
}
