package tasks

import (
	"context"
)

type TaskType string

const (
	RunnableTask TaskType = "FireAndForgetTask"
)

// Runnable task
type Task interface {
	Name() string
	Run(ctx context.Context) error
}

type runnableTask struct {
	name   string
	action TaskAction
}

func NewTask(name string, action TaskAction) Task {
	return &runnableTask{
		name:   name,
		action: action,
	}
}

func (t *runnableTask) Name() string {
	return t.name
}

func (t *runnableTask) Run(ctx context.Context) error {
	err := t.action(ctx)
	return err
}
