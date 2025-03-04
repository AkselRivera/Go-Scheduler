package scheduler

import (
	"container/heap"
	"errors"
	"go-scheduler/models"
	"sync"
	"time"

	"github.com/gofiber/fiber/v2/log"
)

var TaskChannel = make(chan *models.Task)
var CancelChannel = make(chan string)
var taskHeap = &models.TaskHeap{}
var mu sync.Mutex
var taskID string
var taskMap = make(map[string]*models.Task)

func AddTask(t *models.Task) (string, error) {
	mu.Lock()
	defer mu.Unlock()

	if _, exists := taskMap[t.ID]; exists {
		return "", errors.New("task with the same ID already exists")
	}

	t.ID = taskID
	taskMap[t.ID] = t

	return t.ID, nil
}

func CancelTask(id string) {
	CancelChannel <- id
}

func TaskRunner() {
	heap.Init(taskHeap)
	for {
		mu.Lock()
		now := time.Now().UTC()

		for taskHeap.Len() > 0 {
			task := taskHeap.Peek()
			if task.ExecutionTime.After(now) {
				break
			}
			heap.Pop(taskHeap)
			go task.Action()
			delete(taskMap, task.ID)
		}

		var sleepDuration time.Duration
		if nextTask := taskHeap.Peek(); nextTask != nil {
			sleepDuration = nextTask.ExecutionTime.Sub(now)
		} else {
			sleepDuration = time.Hour
		}
		mu.Unlock()

		select {
		case newTask := <-TaskChannel:
			mu.Lock()
			heap.Push(taskHeap, newTask)
			mu.Unlock()
		case taskID := <-CancelChannel:
			mu.Lock()
			if task, exists := taskMap[taskID]; exists && task.Index >= 0 {
				heap.Remove(taskHeap, task.Index)
				delete(taskMap, taskID)
			} else {
				log.Warn("Task with ID %s does not exist or has invalid index\n", taskID)
			}
			mu.Unlock()
		case <-time.After(sleepDuration):
		}
	}
}

func GetAllTasks() []*models.Task {
	mu.Lock()
	defer mu.Unlock()

	tasks := make([]*models.Task, len(*taskHeap))
	copy(tasks, *taskHeap)
	return tasks
}
