package entity

import "time"

const (
	StatusDone    = "done"
	StatusNotDone = "not done"
)

type Task struct {
	ID          int       `json:"id,omitempty"`
	UserID      int       `json:"user_id,omitempty"`
	Title       string    `json:"title,omitempty"`
	Description string    `json:"description,omitempty"`
	Date        time.Time `json:"date,omitempty"`
	Status      string    `json:"status,omitempty"`
}
