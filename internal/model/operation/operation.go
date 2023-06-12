package operation

import (
	"context"
	"fmt"

	"github.com/microsoft/commercial-marketplace-offer-deploy/internal/model"
)

// a executable operation with an execution context
type OperationFunc func(context ExecutionContext) error

// remarks: Invoked Operation decorator+visitor
type Operation struct {
	model.InvokedOperation
	manager *OperationManager
	task    OperationTask
}

func (o *Operation) Context() context.Context {
	return o.manager.ctx
}

func (o *Operation) Running() error {
	o.manager.log.Info("Marking operation as running")

	o.InvokedOperation.Running()
	err := o.manager.saveChanges(true)
	if err != nil {
		o.manager.log.Errorf("failed to save running changes for operation: %v", err)
	}

	return nil
}

func (o *Operation) Complete() error {
	o.manager.log.Info("Marking operation as complete")

	o.InvokedOperation.Complete()
	return o.manager.saveChanges(true)
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
	o.manager.log.Info("Marking operation as failed")

	o.InvokedOperation.Failed()
	return o.saveChangesWithoutNotification()
}

func (o *Operation) Success() error {
	o.manager.log.Info("Marking operation as success")

	o.InvokedOperation.Success()
	return o.saveChangesWithoutNotification()
}

func (o *Operation) Schedule() error {
	o.manager.log.Info("Marking operation as scheduled")

	err := o.InvokedOperation.Schedule()
	if err != nil {
		return err
	}

	err = o.manager.saveChanges(true)
	if err != nil {
		return fmt.Errorf("failed to save schedule changes for operation: %v", err)
	}

	err = o.manager.dispatch()
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
	if o.InvokedOperation.AttemptsExceeded() {
		return o.Complete()
	}
	o.manager.log.Info("Retrying operation")
	return o.Schedule()
}

// provides access to latest instance of associated deployment
func (o *Operation) Deployment() *model.Deployment {
	return o.manager.deployment()
}

// sets the operation's task function
func (o *Operation) Task(task OperationTask) {
	o.task = task
}

// executes the operation
func (o *Operation) Execute() error {
	fn := o.getFunc()
	if fn == nil {
		return fmt.Errorf("invalid runtime operation. the Operation has no task function to execute")
	}

	context := newExecutionContext(o)
	executor := NewExecutor(fn)

	return executor.Execute(context)
}

func (o *Operation) saveChangesWithoutNotification() error {
	return o.manager.saveChanges(false)
}

func (o *Operation) getFunc() OperationFunc {
	if o.IsContinuation() {
		return o.task.Continue
	}
	return o.task.Run
}
