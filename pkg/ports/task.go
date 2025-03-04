package ports

import (
	"go-scheduler/pkg/domain"
)

type TaskService struct {
	TaskHeap      TaskHeapInterface
	TaskMap       map[string]*domain.Task
	TaskChannel   chan *domain.Task
	CancelChannel chan string
}

type TaskServiceInterface interface {
	AddTask(t *domain.Task) (string, error)
	CancelTask(id string)
	GetAllTasks(page, pageSize int) ([]*domain.Task, int)
	TaskRunner()
}
