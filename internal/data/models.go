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
	Status         string `json:"status"`
	DeploymentName string `json:"deploymentName"`
}

type Deployment struct {
	gorm.Model
	Name     string         `json:"name"`
	Status   string         `json:"status"`
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
	Name         string         `json:"name"`
	DeploymentId uint           `json:"deploymentId"`
	Retries      int            `json:"retries"`
	Attempts     int            `json:"attempts"`
	Parameters   map[string]any `json:"parameters" gorm:"json"`
	Result       any            `json:"result" gorm:"json"`
	Status       string         `json:"status"`
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
