package errorz

// APIError представляет структуру ответа Handler'ов для дальнейшей обработки в middleware.
type APIError struct {
	Status int
	Err    error
	Msg    string
}

func (e APIError) Error() string {
	return e.Err.Error()
}
