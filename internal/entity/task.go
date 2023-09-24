package entity

import "time"

const (
	StatusDone    = "done"
	StatusNotDone = "not done"
)

type Task struct {
	UserID      int       `json:"user_id,omitempty"`
	Title       string    `json:"title,omitempty"`
	Description string    `json:"description,omitempty"`
	Data        time.Time `json:"data,omitempty"`
	Status      string    `json:"status,omitempty"`
}
