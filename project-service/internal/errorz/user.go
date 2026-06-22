package errorz

import "errors"

var (
	ErrUserIDNotSet = errors.New("user_id not set")
)
