package domain

import "time"

type Task struct {
	ID            string    `json:"job_id"`
	ExecutionTime time.Time `json:"execution_time"`
	Action        func()    `json:"-"`
	Index         int       `json:"-"`
}
