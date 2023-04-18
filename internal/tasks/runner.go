package tasks

import (
	"context"
	"log"
	"sync"
)

type TaskRunner interface {
	Add(task Task)
	Start() error
}

type runner struct {
	tasks []Task
}

func NewTaskRunner() TaskRunner {
	return &runner{}
}

func (r *runner) Add(task Task) {
	r.tasks = append(r.tasks, task)
}

func (r *runner) Start() error {
	taskCount := len(r.tasks)
	waitGroup := sync.WaitGroup{}
	waitGroup.Add(taskCount)

	ctx := context.Background()

	// just execute all tasks immediately in parallel

	for i := 0; i < taskCount; i++ {
		go func(i int) {
			defer waitGroup.Done()
			task := r.tasks[i]
			err := task.Run(ctx)

			if err != nil {
				log.Printf("Error running task %v: %v", task, err)
			}
		}(i)
	}
	waitGroup.Wait()
	return nil
}
