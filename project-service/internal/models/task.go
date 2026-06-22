package models

import (
	"project-service/internal/errorz"
	"time"
)

type Task struct {
	ID          string     `json:"id" db:"id"`
	Title       string     `json:"title" db:"title"`
	Description string     `json:"description" db:"description"`
	Status      TaskStatus `json:"status" db:"status"`
	ProjectId   string     `json:"project_id" db:"project_id"`
	ExecutorId  string     `json:"executor_id" db:"executor_id"`
	Deadline    time.Time  `json:"deadline" db:"deadline"`
	CreatedAt   time.Time  `json:"created_at" db:"created_at"`
	UpdatedAt   time.Time  `json:"updated_at" db:"updated_at"`
}

// TaskStatus представляет тип статуса задачи.
type TaskStatus string

const (
	StatusInBackLog  TaskStatus = "IN_BACKLOG"
	StatusToDo       TaskStatus = "TO_DO"
	StatusInProgress TaskStatus = "IN_PROGRESS"
	StatusDone       TaskStatus = "DONE"
)

// Validate валидирует статус заявки.
func (s TaskStatus) Validate() error {
	switch s {
	case StatusInBackLog, StatusToDo, StatusInProgress, StatusDone:
		return nil
	default:
		return errorz.ErrInvalidTaskStatus
	}
}

const (
	CreateTaskMsg = "%d:-:В проекте %s создали новое задание %s!"
	UpdateTaskMsg = "%d:-:В проекте %s обновили задание!"
	DeleteTaskMsg = "%d:-:В проекте %s удалили задание!"
)
