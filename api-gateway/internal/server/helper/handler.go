package helper

import (
	"errors"
	"net/http"

	"github.com/rs/zerolog/log"

	"api-gateway/internal/models/response"

	"api-gateway/internal/errorz"
)

// MakeHandler преобразует функцию, принимающую http.ResponseWriter и *http.Request и возвращающую ошибку,
// в http.HandlerFunc, который обрабатывает запросы, обрабатывает возможные ошибки и возвращает соответствующий ответ.
func MakeHandler(fn func(http.ResponseWriter, *http.Request) error) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if err := fn(w, r); err != nil {
			var e errorz.APIError
			if errors.As(err, &e) {
				log.Error().Msg(err.Error())
				WriteJSON(w, e.Status, response.Error{Error: e.Msg})
			}
		}
	}
}
