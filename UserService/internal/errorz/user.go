package errorz

import "errors"

var (
	ErrUserNotFound         = errors.New("user not found")
	ErrPasswordDoesNotMatch = errors.New("password does not match")
	ErrUserIDNotSet         = errors.New("user_id not set")
	ErrUserAccessDenied     = errors.New("user access denied")
)

// вот тут применение errors.New разумно, только почему пакет называется errorZ?)))
// =) errors просто занято уже, а пакет с ошибками, которые будут использоваться в проекте, должен как-то об этом говорить,
// поэтому errorZ. А есть ли другие практики, как его называть?
