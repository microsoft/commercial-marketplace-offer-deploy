package operation

import (
	"context"

	"github.com/microsoft/commercial-marketplace-offer-deploy/internal/model"
	log "github.com/sirupsen/logrus"
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

func (c *ExecutionContext) Success() {
	err := c.operation.Success()
	if err != nil {
		log.Errorf("error updating invoked operation to success: %s", err.Error())
	}
}

func (c *ExecutionContext) Complete() {
	c.operation.Complete()
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

func (c *ExecutionContext) Value(v any) error {
	return c.operation.Value(v)
}

func (c *ExecutionContext) Attribute(key model.AttributeKey, v any) error {
	return c.operation.Attribute(key, v)
}

func (c *ExecutionContext) Operation() *Operation {
	return c.operation
}

func newExecutionContext(operation *Operation) *ExecutionContext {
	return &ExecutionContext{
		operation: operation,
	}
}
