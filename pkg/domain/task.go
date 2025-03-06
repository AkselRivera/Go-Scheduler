package domain

import "time"

type Task struct {
	ID            string    `json:"id"`
	ExecutionTime time.Time `json:"execution_time"`
	Action        func()    `json:"action"`
	Index         int       `json:"-"`
}
