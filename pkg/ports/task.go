package ports

import (
	"sync"

	"github.com/AkselRivera/go-scheduler/pkg/domain"
)

type TaskService struct {
	TaskHeap      TaskHeapInterface
	TaskMap       map[string]*domain.Task
	TaskChannel   chan *domain.Task
	CancelChannel chan string

	Mu *sync.Mutex
}

type TaskServiceInterface interface {
	AddTask(t *domain.Task) (string, error)
	CancelTask(id string)
	GetAllTasks(page, pageSize int) ([]*domain.Task, int)
	TaskRunner()
}
