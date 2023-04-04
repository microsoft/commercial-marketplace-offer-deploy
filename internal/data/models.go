package data

import (
	"strconv"
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
	base.ID = uuid.New()
	return nil
}

type Stage struct {
	gorm.Model
	Name string `json:"name"`
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

// Gets the azure deployment name suitable for azure deployment
// format - modm.<deploymentId>-<deploymentName>
func (d *Deployment) GetAzureDeploymentName() string {
	return "modm." + strconv.FormatUint(uint64(d.ID), 10) + "-" + d.Name
}

type EventSubscription struct {
	BaseWithGuidPrimaryKey
	Callback  string `json:"callback"`
	Name      string `json:"name" gorm:"unique"`
	EventType string `json:"eventType"`
	ApiKey    string `json:"authKey"`
}

type InvokedOperation struct {
	//TODO: stub out for fetching in the operator
}
