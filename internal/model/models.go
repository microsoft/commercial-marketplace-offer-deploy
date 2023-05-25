package model

import (
	"time"

	"github.com/google/uuid"
	"github.com/microsoft/commercial-marketplace-offer-deploy/sdk"
	"gorm.io/gorm"
)

// Base contains common columns for all tables.
type BaseWithGuidPrimaryKey struct {
	ID        uuid.UUID `gorm:"type:uuid;primary_key;"`
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt *time.Time `sql:"index"`
}

// BeforeCreate will set a UUID rather than numeric ID.
func (base *BaseWithGuidPrimaryKey) BeforeCreate(tx *gorm.DB) error {
	if base.ID == uuid.Nil {
		base.ID = uuid.New()
	}
	return nil
}

type Stage struct {
	BaseWithGuidPrimaryKey
	Name           string `json:"name"`
	DeploymentName string `json:"deploymentName"`

	// the default number of retries
	Retries uint `json:"retries"`
}

type Deployment struct {
	gorm.Model
	Name     string         `json:"name"`
	Template map[string]any `json:"template" gorm:"json"`
	Stages   []Stage        `json:"stages" gorm:"json"`

	// azure properties
	SubscriptionId string `json:"subscriptionId"`
	ResourceGroup  string `json:"resourceGroup"`
	Location       string `json:"location"`
}

type EventHook struct {
	BaseWithGuidPrimaryKey
	Callback string `json:"callback"`
	Name     string `json:"name" gorm:"unique"`
	ApiKey   string `json:"authKey"`
}

type InvokedOperation struct {
	BaseWithGuidPrimaryKey
	Name         string `json:"name"`
	DeploymentId uint   `json:"deploymentId"`
	// the correlation id used to track the operation (the correlation id will be set by default to the value on the azure deployment)
	CorrelationId *uuid.UUID              `json:"correlationId" gorm:"type:uuid"`
	Retries       int                     `json:"retries"`
	Attempts      int                     `json:"attempts"`
	Parameters    map[string]any          `json:"parameters" gorm:"json"`
	Result        any                     `json:"result" gorm:"json"`
	Status        string                  `json:"status"`
	Errors        []InvokedOperationError `json:"errors" gorm:"json"`
}

type InvokedOperationError struct {
	Error      string    `json:"error"`
	OccurredAt time.Time `json:"occurredAt"`
	Attempt    int       `json:"attempt"`
}

// appends an error to the list of errors
func (o *InvokedOperation) Error(err error) {
	o.Errors = append(o.Errors, InvokedOperationError{
		Error:      err.Error(),
		OccurredAt: time.Now().UTC(),
		Attempt:    o.Attempts,
	})
}

// increment the number of attempts and set the status to running
func (o *InvokedOperation) Running() {
	o.Attempts++
	o.Status = sdk.StatusRunning.String()
}

func (o *InvokedOperation) BeforeCreate(tx *gorm.DB) error {
	if o.Result == nil {
		o.Result = ""
	}
	err := o.BaseWithGuidPrimaryKey.BeforeCreate(tx)

	if err != nil {
		return err
	}
	return nil
}

func (o *InvokedOperation) BeforeUpdate(tx *gorm.DB) error {
	if o.Result == nil {
		o.Result = ""
	}
	return nil
}
