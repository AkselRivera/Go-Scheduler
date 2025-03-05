package infrastructure

import (
	"sync"

	"github.com/AkselRivera/go-scheduler/pkg/services/task"
	"github.com/AkselRivera/go-scheduler/pkg/services/task_heap"
)

func StartScheduler() *task.TaskService {
	mu := &sync.Mutex{}

	taskHeap := &task_heap.TaskHeap{}
	taskScheduler := task.NewTaskService(taskHeap, mu)

	var wg sync.WaitGroup
	wg.Add(1)

	go taskScheduler.TaskRunner(&wg)

	wg.Wait()

	return taskScheduler
}
