package errorz

import (
	"fmt"

	"google.golang.org/grpc/codes"
)

var (
	ErrBadRequest = fmt.Errorf("bad request")
)

func BadRequest() *Error {
	return &Error{err: ErrBadRequest, statusCode: int(codes.InvalidArgument), from: "grpc"}
}
