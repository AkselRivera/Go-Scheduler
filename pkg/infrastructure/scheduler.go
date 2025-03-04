package infrastructure

import (
	"go-scheduler/pkg/services/task"
	"go-scheduler/pkg/services/task_heap"
)

func StartScheduler() *task.TaskService {
	// Iniciar el servicio de scheduler
	taskHeap := &task_heap.TaskHeap{}
	taskScheduler := task.NewTaskService(taskHeap)
	go taskScheduler.TaskRunner()

	return taskScheduler
}
