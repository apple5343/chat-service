package errorz

import "errors"

var (
	ErrProjectNotFound        = errors.New("project not found")
	ErrProjectAccessForbidden = errors.New("project access forbidden")
)
