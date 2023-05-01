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
	Name   string `json:"name"`
	Status string `json:"status"`
}

type Deployment struct {
	gorm.Model
	Name           string         `json:"name"`
	SubscriptionId string         `json:"subscriptionId"`
	ResourceGroup  string         `json:"resourceGroup"`
	Location       string         `json:"location"`
	Status         string         `json:"status"`
	Template       map[string]any `json:"template" gorm:"json"`
	Stages         []Stage        `json:"stages" gorm:"json"`
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
	Parameters   map[string]any `json:"parameters" gorm:"json"`
	Result       any            `json:"result" gorm:"json"`
	Status       string         `json:"status"`
}
