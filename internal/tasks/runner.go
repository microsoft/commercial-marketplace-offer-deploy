package tasks

import (
	"context"
	"sync"

	log "github.com/sirupsen/logrus"
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
			defer recoverPanic()
			task := r.tasks[i]
			log.Printf("Running task %s", task.Name())
			err := task.Run(ctx)

			if err != nil {
				log.Printf("Error running task %v: %v", task.Name(), err)
			}
		}(i)
	}
	waitGroup.Wait()
	return nil
}

// create a function that catches panics and logs them
func recoverPanic() {
	if r := recover(); r != nil {
		log.Printf("Recovered from panic: %v", r)
	}
}
