package task_heap

import "go-scheduler/pkg/domain"

func (th *TaskHeap) Push(x interface{}) {
	n := len(*th)
	task := x.(*domain.Task)
	task.Index = n
	*th = append(*th, task)
}
