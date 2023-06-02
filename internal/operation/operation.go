package operation

import (
	"context"
	"fmt"

	"github.com/microsoft/commercial-marketplace-offer-deploy/internal/model"
)

// a executable operation with an execution context
type OperationFunc func(context *ExecutionContext) error

// remarks: Invoked Operation decorator+visitor
type Operation struct {
	model.InvokedOperation
	service *operationService
	do      OperationFunc
}

func (o *Operation) Context() context.Context {
	return o.service.ctx
}

func (o *Operation) Running() error {
	changed := o.InvokedOperation.Running()
	if changed {
		return o.service.saveChanges(true)
	}
	return nil
}

func (o *Operation) Attribute(key model.AttributeKey, v any) error {
	o.InvokedOperation.Attribute(key, v)
	return o.saveChangesWithoutNotification()
}

func (o *Operation) Value(v any) error {
	o.InvokedOperation.Value(v)
	return o.saveChangesWithoutNotification()
}

func (o *Operation) Failed() error {
	o.InvokedOperation.Failed()
	return o.service.saveChanges(true)
}

func (o *Operation) Success() error {
	o.InvokedOperation.Success()
	return o.service.saveChanges(true)
}

func (o *Operation) Schedule() error {
	err := o.InvokedOperation.Schedule()
	if err != nil {
		return err
	}

	err = o.service.saveChanges(true)
	if err != nil {
		return fmt.Errorf("failed to save schedule changes for operation: %v", err)
	}

	err = o.service.dispatch()
	if err != nil {
		return fmt.Errorf("failed to schedule operation: %v", err)
	}

	return nil
}

func (o *Operation) SaveChanges() error {
	return o.saveChangesWithoutNotification()
}

// Attempts to trigger a retry of the operation, if the operation has a retriable state
func (o *Operation) Retry() error {
	if !o.InvokedOperation.AttemptsExceeded() {
		return nil
	}

	return o.Schedule()
}

// provides access to latest instance of associated deployment
func (o *Operation) Deployment() *model.Deployment {
	return o.service.deployment()
}

// sets the operation's execution function
func (o *Operation) Do(fn OperationFunc) {
	o.do = fn
}

// executes the operation
func (o *Operation) Execute() error {
	context := newExecutionContext(o)
	executor := NewExecutor(o.do)

	return executor.Execute(context)
}

func (o *Operation) saveChangesWithoutNotification() error {
	return o.service.saveChanges(false)
}
