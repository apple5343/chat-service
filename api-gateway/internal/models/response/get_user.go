package response

import (
	"time"
)

type GetUser struct {
	ID        string    `json:"id" db:"id"`
	Email     string    `json:"email" db:"email"`
	Grade     string    `json:"grade" db:"grade"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
}
