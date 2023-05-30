package operation

import (
	"context"
)

// Context object for interacting with an operation execution
//
//	remarks: flyweight of an Operation
type ExecutionContext struct {
	operation *Operation
}

func (c *ExecutionContext) Context() context.Context {
	return c.operation.Context()
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

func newExecutionContext(operation *Operation) *ExecutionContext {
	return &ExecutionContext{
		operation: operation,
	}
}
