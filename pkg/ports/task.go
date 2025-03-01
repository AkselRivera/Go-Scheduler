package ports

import (
	"errors"
	"go-scheduler/pkg/domain"
	"time"
)

type TaskService struct {
	taskHeap      TaskHeapInterface
	taskMap       map[string]*domain.Task
	TaskChannel   chan *domain.Task
	CancelChannel chan string
}

func NewTaskService(taskHeap TaskHeapInterface) *TaskService {
	return &TaskService{
		taskHeap:      taskHeap,
		taskMap:       make(map[string]*domain.Task),
		TaskChannel:   make(chan *domain.Task),
		CancelChannel: make(chan string),
	}
}

func (ts *TaskService) AddTask(t *domain.Task) (string, error) {
	if _, exists := ts.taskMap[t.ID]; exists {
		return "", errors.New("task with the same ID already exists")
	}

	ts.taskMap[t.ID] = t
	ts.TaskChannel <- t // Enviar la tarea al canal para ser procesada
	return t.ID, nil
}

func (ts *TaskService) CancelTask(id string) {
	ts.CancelChannel <- id
}

func (ts *TaskService) GetAllTasks() []*domain.Task {
	tasks := make([]*domain.Task, len(ts.taskMap))
	i := 0
	for _, task := range ts.taskMap {
		tasks[i] = task
		i++
	}
	return tasks
}

func (ts *TaskService) TaskRunner() {
	for {
		now := time.Now().UTC()

		// Ejecutar tareas del heap
		for ts.taskHeap.Len() > 0 {
			task := ts.taskHeap.Peek()
			if task.ExecutionTime.After(now) {
				break
			}
			// Ejecutar la tarea
			task.Action()
			delete(ts.taskMap, task.ID)
			ts.taskHeap.Pop()
		}

		var sleepDuration time.Duration
		if nextTask := ts.taskHeap.Peek(); nextTask != nil {
			sleepDuration = nextTask.ExecutionTime.Sub(now)
		} else {
			sleepDuration = time.Hour
		}

		// Escuchar los canales
		select {
		case newTask := <-ts.TaskChannel:
			ts.taskHeap.Push(newTask)
		case taskID := <-ts.CancelChannel:
			if task, exists := ts.taskMap[taskID]; exists && task.Index >= 0 {
				// Remover tarea del heap
				ts.taskHeap.Pop()
				delete(ts.taskMap, taskID)
			}
		case <-time.After(sleepDuration):
		}
	}
}
