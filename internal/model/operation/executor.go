package operation

import (
	"context"
	"fmt"
	"time"

	"github.com/microsoft/commercial-marketplace-offer-deploy/pkg/deployment"
	"github.com/microsoft/commercial-marketplace-offer-deploy/sdk"
	log "github.com/sirupsen/logrus"
)

// RetriableError is a custom error that contains a positive duration for the next retry
type RetriableError struct {
	Err        error
	RetryAfter time.Duration
}

// Error returns error message and a Retry-After duration
func (e *RetriableError) Error() string {
	return fmt.Sprintf("%s (retry after %v)", e.Err.Error(), e.RetryAfter)
}

// this is so the dry run can be tested, detaching actual dry run implementation
type DryRunFunc func(context.Context, *deployment.AzureDeployment) (*sdk.DryRunResult, error)

// Executor is the interface for the actual execution of a logically invoked operation from the API
// Requestor --> invoke this operation --> enqueue --> executor --> execute the operation
type Executor interface {
	Execute(context ExecutionContext) error
}

// default implementation of an operation executor
type executor struct {
	operation OperationFunc
}

// default implementation for executing an operation
func (exe *executor) Execute(executionContext ExecutionContext) error {
	if reasons, ok := executionContext.Operation().IsExecutable(); !ok {
		log.Infof("Operation is not in an executable state: %s", reasons)
		return nil
	}

	err := executionContext.Running()

	if err != nil {
		log.Errorf("error updating operation to running: %s", err.Error())
		return err
	}

	do := WithLogging(exe.operation)
	err = do(executionContext)

	if err != nil {
		executionContext.Error(err)

		if executionContext.Operation().AttemptsExceeded() {
			err = executionContext.Failed()
			executionContext.Complete()
		} else {
			retryErr := executionContext.Retry()
			if retryErr != nil {
				log.Errorf("attempt to retry operation caused error: %s", retryErr.Error())
			}
		}

		return err
	}

	executionContext.Success()
	executionContext.Complete()

	return nil
}

// default operations executor that executions the operation(s) in sequence with logging and default retry logic
//
//	remarks: if any of the operations return an error, the executor considers this a failure and will not execute
func NewExecutor(operation OperationFunc) Executor {
	return &executor{
		operation: operation,
	}
}

func WithLogging(operation OperationFunc) OperationFunc {
	return func(context ExecutionContext) error {
		logger := log.WithFields(
			log.Fields{
				"operation": fmt.Sprintf("%+v", context.Operation()),
			})
		logger.Debug("Executing operation")
		err := operation(context)
		logger.Debug("Execution of operation done")

		if err != nil {
			logger.WithError(err).Error("execution failed")
		}
		return err
	}
}
