package operation

import (
	"context"
)

type ExecutionContext struct {
	ctx       context.Context
	operation *Operation
}

func (c *ExecutionContext) Context() context.Context {
	return c.ctx
}

func (c *ExecutionContext) Running() error {
	return c.operation.Running()
}

func (c *ExecutionContext) Error(err error) {
	c.operation.Error(err)
}

func (c *ExecutionContext) Success() error {
	return c.operation.Success()
}

func (c *ExecutionContext) SaveChanges() error {
	return c.operation.SaveChanges()
}

func (c *ExecutionContext) Failed() error {
	return c.operation.Failed()
}

func (c *ExecutionContext) Retry() error {
	return c.operation.Retry()
}

func (c *ExecutionContext) Value(v any) {
	c.operation.Value(v)
}

func (c *ExecutionContext) Operation() *Operation {
	return c.operation
}

func NewExecutionContext(ctx context.Context, operation *Operation) *ExecutionContext {
	return &ExecutionContext{
		ctx:       ctx,
		operation: operation,
	}
}
