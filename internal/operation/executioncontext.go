package operation

import (
	"context"
)

type ExecutionContext struct {
	ctx              context.Context
	invokedOperation *Operation
}

func (c *ExecutionContext) Context() context.Context {
	return c.ctx
}

func (c *ExecutionContext) Running() error {
	return c.invokedOperation.Running()
}

func (c *ExecutionContext) Error(err error) {
	c.invokedOperation.Error(err)
}

func (c *ExecutionContext) Success() error {
	return c.invokedOperation.Success()
}

func (c *ExecutionContext) Failed() error {
	return c.invokedOperation.Failed()
}

func (c *ExecutionContext) Retry() error {
	return c.invokedOperation.Retry()
}

func (c *ExecutionContext) Value(v any) {
	c.invokedOperation.Value(v)
}

func (c *ExecutionContext) InvokedOperation() *Operation {
	return c.invokedOperation
}

func NewExecutionContext(ctx context.Context, invokedOperation *Operation) *ExecutionContext {
	return &ExecutionContext{
		ctx:              ctx,
		invokedOperation: invokedOperation,
	}
}
