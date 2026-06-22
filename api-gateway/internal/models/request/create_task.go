package request

import "time"

type CreateTask struct {
	Title       string    `json:"title"`
	Description string    `json:"description"`
	Status      string    `json:"status"`
	Executor    string    `json:"executor"`
	Deadline    time.Time `json:"deadline"`
}
