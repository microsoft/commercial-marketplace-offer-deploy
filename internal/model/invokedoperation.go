package model

import (
	"fmt"
	"time"

	"github.com/microsoft/commercial-marketplace-offer-deploy/sdk"
)

const DefaultNumberOfRetries = 3

// InvokedOperation is a record of an operation that was invoked by the operator
//
//	remarks:
//		attributes - specify information about or that control the operation's use and behavior
type InvokedOperation struct {
	BaseWithGuidPrimaryKey
	Name         string                           `json:"name"`
	DeploymentId uint                             `json:"deploymentId"`
	Attributes   []InvokedOperationAttribute      `json:"attributes"`
	Retries      uint                             `json:"retries"`
	Attempts     uint                             `json:"attempts"`
	Parameters   map[string]any                   `json:"parameters" gorm:"json"`
	Results      map[uint]*InvokedOperationResult `json:"results" gorm:"json"`

	// the current or final status of the operation
	Status string `json:"status"`
}

type InvokedOperationResult struct {
	Attempt     uint      `json:"attempt"`
	Error       string    `json:"error"`
	Value       any       `json:"value" gorm:"json"`
	StartedAt   time.Time `json:"startedAt"`
	CompletedAt time.Time `json:"completedAt"`
	Status      string    `json:"status"`
}

// can the operation be executed? if not, then reasons are returned
func (o *InvokedOperation) IsExecutable() ([]string, bool) {
	reasons := []string{}

	isRunning := o.IsRunning()
	if isRunning {
		reasons = append(reasons, "operation is already running")
	}
	attemptsExceeded := o.AttemptsExceeded()
	if attemptsExceeded {
		reasons = append(reasons, fmt.Sprintf("operation has exceeded the maximum number of attempts (%d)", o.Retries))
	}
	return reasons, (!isRunning && !attemptsExceeded)
}

func (o *InvokedOperation) IsRunning() bool {
	return o.Status == sdk.StatusRunning.String()
}

// increment the number of attempts and set the status to running
func (o *InvokedOperation) Running() (error, bool) {
	if o.IsRunning() { //already running, so do nothing
		return nil, true
	}

	o.incrementAttempts()

	if o.AttemptsExceeded() {
		return fmt.Errorf("cannot run operation, %d of %d attemps reached", +o.Attempts, o.Retries), false
	}

	o.Status = sdk.StatusRunning.String()
	o.appendResult()

	return nil, true
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
	return o.Attempts >= o.Retries
}

func (o *InvokedOperation) IsRetry() bool {
	return o.Attempts > 1
}

func (io *InvokedOperation) IsRetriable() bool {
	return !io.IsRunning() && !io.AttemptsExceeded()
}

// sets the status to failed for the operation and the latest attempt's result
func (o *InvokedOperation) Failed() {
	o.setStatus(sdk.StatusFailed.String())
}

// sets the status to success for the operation and the latest attempt's result
func (o *InvokedOperation) Success() {
	o.setStatus(sdk.StatusSuccess.String())
}

func (o *InvokedOperation) Schedule() error {
	if o.AttemptsExceeded() {
		return fmt.Errorf("cannot schedule operation, %d of %d attemps reached", o.Attempts, o.Retries)
	}
	o.setStatus(sdk.StatusScheduled.String())
	return nil
}

func (o *InvokedOperation) setStatus(status string) {
	o.Status = status
	result := o.LatestResult()
	result.Status = status
	result.CompletedAt = time.Now().UTC()
}

func (o *InvokedOperation) incrementAttempts() {
	if o.AttemptsExceeded() {
		return
	}
	o.Attempts++
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
