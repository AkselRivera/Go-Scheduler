package main

import (
	"fmt"
	"sync"
	"time"

	"github.com/AkselRivera/go-scheduler/pkg/domain"
	"github.com/AkselRivera/go-scheduler/pkg/infrastructure"
)

func main() {
	taskScheduler := infrastructure.StartScheduler()
	taskScheduler.AddTask(&domain.Task{ID: "0", ExecutionTime: time.Now(), Action: func() { fmt.Println("Task 0 executed", time.Now()) }})
	tt, er := taskScheduler.AddTask(&domain.Task{ID: "01", ExecutionTime: time.Now().Add(time.Second * -10), Action: func() { fmt.Println("Task 01 executed", time.Now()) }})

	if er != nil {
		fmt.Println("Error adding task", er)
	} else {
		fmt.Println("Task 01 added at", time.Now(), tt)
	}

	wg := sync.WaitGroup{}

	t, total := taskScheduler.GetAllTasks(1, 25)
	fmt.Println("Getting all tasks", t, "total", total)

	wg.Add(1)
	s, err := taskScheduler.AddTask(&domain.Task{ID: "1", ExecutionTime: time.Now().Add(time.Minute * 1), Action: func() {
		fmt.Println("Task 1 executed", time.Now())
		wg.Done()
	}})
	if err != nil {
		panic(err)
	}

	fmt.Println("Task 1 added at", time.Now(), s)

	wg.Add(1)
	s, err = taskScheduler.AddTask(&domain.Task{ID: "2", ExecutionTime: time.Now().Add(time.Second * 30), Action: func() {
		fmt.Println("Task 2 executed ", time.Now())
		wg.Done()

	}})
	if err != nil {
		panic(err)
	}
	fmt.Println("Task 2 added at", time.Now(), s)

	funcWithParams := func(x, y string, wg *sync.WaitGroup) {
		fmt.Println("Task 3 executed ", time.Now(), x, y)
		wg.Done()
	}

	wg.Add(1)
	s, err = taskScheduler.AddTask(&domain.Task{ID: "3", ExecutionTime: time.Now().Add(time.Second * 10),
		Action: func() {
			funcWithParams("hello", "world", &wg)
		}})
	if err != nil {
		panic(err)
	}

	fmt.Println("Task 3 added at", time.Now(), s)

	fmt.Println("Getting all tasks")
	getAllTasks(taskScheduler.GetAllTasks(1, 25))

	wg.Wait()

	fmt.Println("Getting all tasks finished")
	getAllTasks(taskScheduler.GetAllTasks(1, 25))

}

func getAllTasks(tasks []*domain.Task, total int) {
	for _, t := range tasks {
		fmt.Println(t, "total:", total)
	}
}
