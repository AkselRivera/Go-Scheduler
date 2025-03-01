package infrastructure

import (
	"go-scheduler/pkg/services/task_heap"
)

func StartScheduler() {
	// Iniciar el servicio de scheduler
	taskHeap := &task_heap.TaskHeap{}
	taskService := application.NewTaskService(taskHeap)
	go taskService.TaskRunner()
}
