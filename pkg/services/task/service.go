package task

import (
	"sync"

	"github.com/AkselRivera/go-scheduler/pkg/domain"
	"github.com/AkselRivera/go-scheduler/pkg/ports"
)

type TaskService ports.TaskService

func NewTaskService(taskHeap ports.TaskHeapInterface, mu *sync.Mutex) *TaskService {
	return &TaskService{
		TaskHeap:      taskHeap,
		TaskMap:       make(map[string]*domain.Task),
		TaskChannel:   make(chan *domain.Task),
		CancelChannel: make(chan string),
		Mu:            mu,
	}
}
