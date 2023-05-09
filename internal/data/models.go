package data

import (
	"time"

	"github.com/google/uuid"
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
	CorrelationId *uuid.UUID     `json:"correlationId" gorm:"type:uuid"`
	Retries       int            `json:"retries"`
	Attempts      int            `json:"attempts"`
	Parameters    map[string]any `json:"parameters" gorm:"json"`
	Result        any            `json:"result" gorm:"json"`
	Status        string         `json:"status"`
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
