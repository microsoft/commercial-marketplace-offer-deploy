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

type EventSubscription struct {
	BaseWithGuidPrimaryKey
	Callback *string `json:"callback"`
	Name     *string `json:"name"`
	Topic    *string `json:"topic"`
}
