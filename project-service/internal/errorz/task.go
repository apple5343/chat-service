package errorz

import "errors"

var (
	ErrTaskNotFound        = errors.New("task not found")
	ErrTaskAccessForbidden = errors.New("task access forbidden")
	ErrInvalidTaskStatus   = errors.New("invalid task status")
)
