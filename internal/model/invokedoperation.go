package model

import (
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/microsoft/commercial-marketplace-offer-deploy/sdk"
)

const DefaultNumberOfRetries = 3

// InvokedOperation is a record of an operation that was invoked by the operator
//
//	remarks:
//		attributes - specify information about or that control the operation's use and behavior
type InvokedOperation struct {
	BaseWithGuidPrimaryKey
	ParentID     *uuid.UUID                       `json:"parentId" gorm:"type:uuid"`
	Name         string                           `json:"name"`
	DeploymentId uint                             `json:"deploymentId"`
	Attributes   []InvokedOperationAttribute      `json:"attributes"`
	Retries      uint                             `json:"retries"`
	Attempts     uint                             `json:"attempts"`
	Parameters   map[string]any                   `json:"parameters" gorm:"json"`
	Results      map[uint]*InvokedOperationResult `json:"results" gorm:"json"`
	Status       string                           `json:"status"` // the current or final status of the operation
	Completed    bool                             `json:"completed"`
}

type InvokedOperationResult struct {
	Attempt     uint      `json:"attempt"`
	Error       string    `json:"error"`
	Value       any       `json:"value" gorm:"json"`
	Status      string    `json:"status"`
	StartedAt   time.Time `json:"startedAt"`
	CompletedAt time.Time `json:"completedAt"`
}

// can the operation be executed? if not, then reasons are returned
func (o *InvokedOperation) IsExecutable() ([]string, bool) {
	reasons := []string{}

	isCompleted := o.IsCompleted()
	if isCompleted {
		reasons = append(reasons, "operation already completed on [%s]", o.LatestResult().CompletedAt.String())
	}

	isRunning := o.IsRunning()
	if isRunning {
		reasons = append(reasons, "operation is already running")
	}

	attemptsExceeded := o.AttemptsExceeded()
	if attemptsExceeded {
		reasons = append(reasons, fmt.Sprintf("operation has exceeded the maximum number of attempts (%d)", o.Retries))
	}
	return reasons, (!isCompleted && !isRunning && !attemptsExceeded)
}

func (o *InvokedOperation) IsRetry() bool {
	return !o.IsCompleted() && o.Attempts > 1
}

func (o *InvokedOperation) IsFirstAttempt() bool {
	return !o.IsCompleted() && o.Attempts == 1
}

func (o *InvokedOperation) IsRunning() bool {
	return o.Status == sdk.StatusRunning.String()
}

func (o *InvokedOperation) IsCompleted() bool {
	return o.Completed
}

func (o *InvokedOperation) IsScheduled() bool {
	return o.Status == sdk.StatusScheduled.String()
}

func (o *InvokedOperation) CompletedAt() (*time.Time, error) {
	if !o.IsCompleted() {
		return nil, fmt.Errorf("operation is not completed")
	}
	completedAt := o.LatestResult().CompletedAt
	return &completedAt, nil
}

// increment the number of attempts and set the status to running
// if it's in a state where it can be run.
// return: returns true if status changed to running
func (o *InvokedOperation) Running() {
	if o.IsRunning() || o.AttemptsExceeded() { //already running, so do nothing
		return
	}

	o.incrementAttempts()
	o.setStatus(sdk.StatusRunning.String())
}

func (o *InvokedOperation) FirstResult() *InvokedOperationResult {
	if len(o.Results) == 0 {
		return o.appendResult()
	}
	if _, ok := o.Results[1]; !ok {
		return o.appendResult()
	}
	return o.Results[1]
}

func (o *InvokedOperation) LatestResult() *InvokedOperationResult {
	if len(o.Results) == 0 {
		return o.appendResult()
	}
	if _, ok := o.Results[o.Attempts]; !ok {
		return o.appendResult()
	}
	return o.Results[o.Attempts]
}

func (o *InvokedOperation) Error(err error) {
	o.LatestResult().Error = err.Error()
	o.setStatus(sdk.StatusError.String())
}

func (o *InvokedOperation) Value(v any) {
	o.LatestResult().Value = v
}

func (o *InvokedOperation) Attribute(key AttributeKey, v any) {
	if o.Attributes == nil {
		o.Attributes = []InvokedOperationAttribute{}
	}

	for i, attr := range o.Attributes {
		if attr.Key == string(key) {
			o.Attributes[i].Value = v
			return
		}
	}

	o.Attributes = append(o.Attributes, NewAttribute(key, v))
}

// does the InvokedOperation have an attribute with the specified key?
func (o *InvokedOperation) HasAttribute(key AttributeKey) bool {
	for _, attr := range o.Attributes {
		if attr.Key == string(key) {
			return true
		}
	}
	return false
}

func (o *InvokedOperation) AttributeValue(key AttributeKey) (any, bool) {
	for _, attr := range o.Attributes {
		if attr.Key == string(key) {
			return attr.Value, true
		}
	}
	return nil, false
}

func (o *InvokedOperation) AttemptsExceeded() bool {
	return o.Attempts > o.Retries
}

func (o *InvokedOperation) IsRetrying() bool {
	return o.Status == string(sdk.StatusScheduled) && (!o.AttemptsExceeded() && o.Attempts > 1)
}

// sets the status to failed for the operation and the latest attempt's result
func (o *InvokedOperation) Failed() error {
	return o.setStatus(sdk.StatusFailed.String())
}

// sets the status to success for the operation and the latest attempt's result
func (o *InvokedOperation) Success() error {
	return o.setStatus(sdk.StatusSuccess.String())
}

func (o *InvokedOperation) Complete() {
	o.Completed = true
	o.LatestResult().CompletedAt = time.Now().UTC()
}

func (o *InvokedOperation) Schedule() error {
	if o.AttemptsExceeded() {
		return fmt.Errorf("cannot schedule operation, %d of %d attemps reached", o.Attempts, o.Retries)
	}
	return o.setStatus(sdk.StatusScheduled.String())
}

func (o *InvokedOperation) setStatus(status string) error {
	// if the operation is complete, the status cannot be set
	if o.IsCompleted() {
		return fmt.Errorf("cannot set status to %s, operation is already complete", status)
	}

	o.Status = status // track the latest status

	//anything but schedule, update the results
	if status != string(sdk.StatusScheduled) {
		result := o.LatestResult()
		result.Status = status
		result.CompletedAt = time.Now().UTC()
	}
	return nil
}

func (o *InvokedOperation) incrementAttempts() {
	if o.AttemptsExceeded() || o.IsCompleted() {
		return
	}
	o.Attempts++
}

func (o *InvokedOperation) CorrelationId() (*uuid.UUID, error) {
	value, ok := o.AttributeValue(AttributeKeyCorrelationId)
	if !ok {
		return nil, errors.New("no correlation id found for operation")
	}
	correlationIdString := fmt.Sprintf("%v", value)
	correlationId, err := uuid.Parse(correlationIdString)
	if err != nil {
		return nil, err
	}
	return &correlationId, nil
}

func (o *InvokedOperation) appendResult() *InvokedOperationResult {
	if o.Results == nil {
		o.Results = make(map[uint]*InvokedOperationResult)
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
