package model

import (
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/microsoft/commercial-marketplace-offer-deploy/sdk"
)

// InvokedOperation is a record of an operation that was invoked by the operator
type InvokedOperation struct {
	BaseWithGuidPrimaryKey
	Name         string `json:"name"`
	DeploymentId uint   `json:"deploymentId"`
	// the correlation id used to track the operation (the correlation id will be set by default to the value on the azure deployment)
	CorrelationId *uuid.UUID                      `json:"correlationId" gorm:"type:uuid"`
	Retries       int                             `json:"retries"`
	Attempts      int                             `json:"attempts"`
	Parameters    map[string]any                  `json:"parameters" gorm:"json"`
	Results       map[int]*InvokedOperationResult `json:"results" gorm:"json"`

	// the current or final status of the operation
	Status string `json:"status"`
}

type InvokedOperationResult struct {
	Attempt     int       `json:"attempt"`
	Error       string    `json:"error"`
	Value       any       `json:"value" gorm:"json"`
	StartedAt   time.Time `json:"startedAt"`
	CompletedAt time.Time `json:"occurredAt"`
	Status      string    `json:"status"`
}

// increment the number of attempts and set the status to running
func (o *InvokedOperation) Running() error {
	o.Attempts++

	if o.AttemptsExceeded() {
		return fmt.Errorf("cannot run operation, %d of %d attemps reached", +o.Attempts, o.Retries)
	}

	o.Status = sdk.StatusRunning.String()
	o.appendResult()

	return nil
}

func (o *InvokedOperation) LatestResult() *InvokedOperationResult {
	if len(o.Results) == 0 {
		return o.appendResult()
	}
	return o.Results[o.Attempts]
}

func (o *InvokedOperation) Error(err error) {
	o.LatestResult().Error = err.Error()
}

func (o *InvokedOperation) Value(v any) {
	o.LatestResult().Value = v
}

func (o *InvokedOperation) AttemptsExceeded() bool {
	return o.Attempts > o.Retries
}

func (o *InvokedOperation) IsRetry() bool {
	return o.Attempts > 1
}

// sets the status to failed for the operation and the latest attempt's result
func (o *InvokedOperation) Failed() {
	o.setStatus(sdk.StatusFailed.String())
}

// sets the status to success for the operation and the latest attempt's result
func (o *InvokedOperation) Success() {
	o.setStatus(sdk.StatusSuccess.String())
}

func (o *InvokedOperation) setStatus(status string) {
	o.Status = status
	result := o.LatestResult()
	result.Status = status
	result.CompletedAt = time.Now().UTC()
}

func (o *InvokedOperation) appendResult() *InvokedOperationResult {
	if o.Results == nil {
		o.Results = make(map[int]*InvokedOperationResult)
	}

	if _, exists := o.Results[o.Attempts]; exists {
		return o.Results[o.Attempts]
	}

	result := &InvokedOperationResult{
		Attempt:   o.Attempts,
		StartedAt: time.Now().UTC(),
		Value:     nil,
		Status:    sdk.StatusRunning.String(),
	}
	o.Results[o.Attempts] = result

	return result
}
