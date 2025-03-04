package task

import (
	"go-scheduler/pkg/domain"
	"go-scheduler/pkg/ports"
	"sync"
)

type TaskService ports.TaskService

var mu sync.Mutex

func NewTaskService(taskHeap ports.TaskHeapInterface) *TaskService {
	return &TaskService{
		TaskHeap:      taskHeap,
		TaskMap:       make(map[string]*domain.Task),
		TaskChannel:   make(chan *domain.Task),
		CancelChannel: make(chan string),
	}
}
