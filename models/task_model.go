package models

import "time"

// Estructura de la tarea con un ID único
type Task struct {
	ID            string    `json:"job_id"`
	ExecutionTime time.Time `json:"execution_time"`
	Action        func()    `json:"-"`
	Index         int       `json:"-"`
}

// Estructura del heap de tareas
type TaskHeap []*Task

func (th TaskHeap) Len() int           { return len(th) }
func (th TaskHeap) Less(i, j int) bool { return th[i].ExecutionTime.Before(th[j].ExecutionTime) }
func (th TaskHeap) Swap(i, j int) {
	th[i], th[j] = th[j], th[i]
	th[i].Index = i
	th[j].Index = j
}

func (th *TaskHeap) Push(x interface{}) {
	n := len(*th)
	task := x.(*Task)
	task.Index = n
	*th = append(*th, task)
}

func (th *TaskHeap) Pop() interface{} {
	old := *th
	n := len(old)
	task := old[n-1]
	old[n-1] = nil  // Evitar memory leak
	task.Index = -1 // Por seguridad
	*th = old[0 : n-1]
	return task
}

// Implementar Peek para revisar el primer elemento sin eliminarlo
func (th *TaskHeap) Peek() *Task {
	if th.Len() == 0 {
		return nil
	}
	return (*th)[0]
}
