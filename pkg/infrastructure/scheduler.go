package infrastructure

import (
	"sync"

	"github.com/AkselRivera/go-scheduler/pkg/services/task"
	"github.com/AkselRivera/go-scheduler/pkg/services/task_heap"
)

func StartScheduler(mu *sync.Mutex) *task.TaskService {
	// Iniciar el servicio de scheduler
	taskHeap := &task_heap.TaskHeap{}
	taskScheduler := task.NewTaskService(taskHeap, mu)
	go taskScheduler.TaskRunner()

	return taskScheduler
}
