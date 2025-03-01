package task_heap

import "go-scheduler/pkg/domain"

func (th *TaskHeap) Peek() *domain.Task {
	if th.Len() == 0 {
		return nil
	}
	return (*th)[0]
}
