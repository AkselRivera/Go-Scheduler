package ports

import "github.com/AkselRivera/go-scheduler/pkg/domain"

type TaskHeapInterface interface {
	Len() int
	Less(i, j int) bool
	Swap(i, j int)

	Push(x interface{})
	Pop() interface{}
	Peek() *domain.Task
}
