package task

import (
	"container/heap"
	"errors"
	"fmt"
	"go-scheduler/pkg/domain"
	"time"
)

func (ts *TaskService) AddTask(t *domain.Task) (string, error) {
	mu.Lock()
	defer mu.Unlock()

	if _, exists := ts.TaskMap[t.ID]; exists {
		return "", errors.New("task with the same ID already exists")
	}

	ts.TaskMap[t.ID] = t
	ts.TaskChannel <- t // Enviar la tarea al canal para ser procesada
	return t.ID, nil
}

func (ts *TaskService) CancelTask(id string) {
	ts.CancelChannel <- id
}

func (ts *TaskService) GetAllTasks() []*domain.Task {
	mu.Lock()
	defer mu.Unlock()

	tasks := make([]*domain.Task, len(ts.TaskMap))
	i := 0
	for _, task := range ts.TaskMap {
		tasks[i] = task
		i++
	}
	return tasks
}

func (ts *TaskService) TaskRunner() {
	heap.Init(ts.TaskHeap)
	for {
		mu.Lock()
		now := time.Now().UTC()

		for ts.TaskHeap.Len() > 0 {
			task := ts.TaskHeap.Peek()
			if task.ExecutionTime.After(now) {
				break
			}

			heap.Pop(ts.TaskHeap)

			go task.Action()
			delete(ts.TaskMap, task.ID)
		}

		var sleepDuration time.Duration
		if nextTask := ts.TaskHeap.Peek(); nextTask != nil {
			fmt.Println("Next task:", nextTask.ID, "at", nextTask.ExecutionTime)
			sleepDuration = nextTask.ExecutionTime.Sub(now)
		} else {
			sleepDuration = time.Hour
		}
		mu.Unlock()

		// Escuchar los canales
		select {
		case newTask := <-ts.TaskChannel:
			mu.Lock()
			heap.Push(ts.TaskHeap, newTask)
			mu.Unlock()
		case taskID := <-ts.CancelChannel:
			mu.Lock()
			if task, exists := ts.TaskMap[taskID]; exists && task.Index >= 0 {
				heap.Remove(ts.TaskHeap, task.Index)
				delete(ts.TaskMap, taskID)
			}
			mu.Unlock()
		case <-time.After(sleepDuration):
		}
	}
}
