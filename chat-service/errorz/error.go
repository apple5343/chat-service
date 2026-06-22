package errorz

import (
	"errors"
	"fmt"
)

var (
	ErrInternalServer = fmt.Errorf("internal server error")
	ErrUnknown        = fmt.Errorf("unknown error")
	ErrForbidden      = fmt.Errorf("forbidden")
)

type Error struct {
	isInternal bool
	from       string
	err        error
	statusCode int
}

func Wrap(err error, from string) *Error {
	return &Error{
		err:  err,
		from: from,
	}
}

func WrapInternal(err error, from string) *Error {
	return &Error{
		isInternal: true,
		err:        err,
		from:       from,
	}
}

func (e *Error) Message() string {
	if e.isInternal {
		return ErrInternalServer.Error()
	}
	return e.err.Error()
}

func (e *Error) StatusCode() int {
	return e.statusCode
}

func (e *Error) FuncName() string {
	return e.from
}

func (e *Error) SetStatusCode(code int) {
	e.statusCode = code
}

func Parse(err error) *Error {
	e, ok := err.(*Error)
	if !ok {
		return nil
	}
	return e
}

func (e *Error) IsInternal() bool {
	return e.isInternal
}

func (e *Error) Unwrap() error {
	return e.err
}

func (e *Error) Error() string {
	return e.err.Error()
}

func (e *Error) Is(target error) bool {
	return e.err == target || (e.err != nil && errors.Is(e.err, target))
}

func (e *Error) As(target any) bool {
	return e.err != nil && errors.As(e.err, target)
}
