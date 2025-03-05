package task

import (
	"github.com/AkselRivera/go-scheduler/pkg/domain"
	"github.com/AkselRivera/go-scheduler/pkg/ports"
	"sync"
)

type TaskService ports.TaskService

func NewTaskService(taskHeap ports.TaskHeapInterface) *TaskService {
	return &TaskService{
		TaskHeap:      taskHeap,
		TaskMap:       make(map[string]*domain.Task),
		TaskChannel:   make(chan *domain.Task),
		CancelChannel: make(chan string),
		Mu:            sync.Mutex{},
	}
}
