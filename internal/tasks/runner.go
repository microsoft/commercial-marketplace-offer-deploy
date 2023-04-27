package tasks

import (
	"context"
	"sync"
	"time"

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
			retryCount := 10 * 60
			retriedExecutor := Retry(executeTask,  retryCount, 10 * time.Second)
			err := retriedExecutor(ctx, task)

			if err != nil {
				log.Printf("task error [%s]: %v", task.Name(), err)
			}
		}(i)
	}
	waitGroup.Wait()
	return nil
}

type Effector func(context.Context, Task) error

func executeTask(ctx context.Context, task Task) error {
	log.Printf("Executing: %s", task.Name())
	err := task.Run(ctx)
	log.Printf("Task Completed: %s", task.Name())
	return err
}

func Retry(effector Effector, retries int, delay time.Duration) Effector {
	return func(ctx context.Context, task Task) error {
		for r:= 0; ; r++ {
			err := effector(ctx, task)
			if err == nil || r >= retries {
				// success or max retries reached
				return err
			}	

			log.Printf("Attempt %d failed, retrying in %v", r + 1, delay)

			select {
				case <- time.After(delay):
				case <- ctx.Done():
					return ctx.Err()
			}
		}
	}
}

// create a function that catches panics and logs them
func recoverPanic() {
	if r := recover(); r != nil {
		log.Printf("Recovered from panic: %v", r)
	}
}
