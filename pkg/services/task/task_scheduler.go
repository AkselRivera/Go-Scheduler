package task

import (
	"container/heap"
	"errors"
	"fmt"
	"sync"
	"time"

	"github.com/AkselRivera/go-scheduler/pkg/domain"
)

func (ts *TaskService) AddTask(t *domain.Task) (string, error) {
	ts.Mu.Lock()
	defer ts.Mu.Unlock()

	if t.ExecutionTime.Before(time.Now().UTC()) {
		return "", errors.New("task execution time cannot be in the past")
	}

	if _, exists := ts.TaskMap[t.ID]; exists {
		return "", fmt.Errorf("task #%s already exists", t.ID)
	}

	ts.TaskChannel <- t
	return t.ID, nil
}

func (ts *TaskService) CancelTask(id string) {
	ts.CancelChannel <- id
}

func (ts *TaskService) GetAllTasks(page, pageSize int) ([]*domain.Task, int) {
	ts.Mu.Lock()
	defer ts.Mu.Unlock()

	if page < 1 {
		page = 1
	}

	if pageSize < 1 || pageSize > 100 {
		pageSize = 25
	}

	totalTasks := len(ts.TaskMap)
	start := (page - 1) * pageSize
	end := start + pageSize
	if start > len(ts.TaskMap) {
		return nil, 0
	}
	if end > len(ts.TaskMap) {
		end = len(ts.TaskMap)
	}

	tasks := make([]*domain.Task, 0, end-start)
	i := 0
	for _, task := range ts.TaskMap {
		if i >= start && i < end {
			tasks = append(tasks, task)
		}
		i++
	}
	return tasks, totalTasks
}

func (ts *TaskService) TaskRunner(wg *sync.WaitGroup) {
	heap.Init(ts.TaskHeap)
	wg.Done()
	for {
		ts.Mu.Lock()
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
			sleepDuration = nextTask.ExecutionTime.Sub(now)
		} else {
			sleepDuration = time.Hour
		}
		ts.Mu.Unlock()

		select {
		case newTask := <-ts.TaskChannel:
			ts.Mu.Lock()
			heap.Push(ts.TaskHeap, newTask)
			ts.Mu.Unlock()
		case taskID := <-ts.CancelChannel:
			ts.Mu.Lock()
			if task, exists := ts.TaskMap[taskID]; exists && task.Index >= 0 {
				heap.Remove(ts.TaskHeap, task.Index)
				delete(ts.TaskMap, taskID)
			}
			ts.Mu.Unlock()
		case <-time.After(sleepDuration):
		}
	}
}
