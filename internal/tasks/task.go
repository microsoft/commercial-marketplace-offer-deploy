package tasks

import (
	"context"
	"log"
)

type TaskType string

const (
	RunnableTask TaskType = "FireAndForgetTask"
)

// Runnable task
type Task interface {
	Run(ctx context.Context) error
}

type runnableTask struct {
	action TaskAction
}

func NewTask(action TaskAction) Task {
	return &runnableTask{
		action: action,
	}
}

func (t *runnableTask) Run(ctx context.Context) error {
	log.Printf("Running task %v", t)
	err := t.action(ctx)
	if err != nil {
		log.Printf("Error running task %v: %v", t, err)
		return err
	}
	return nil
}
