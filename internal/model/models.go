package model

import (
	"gorm.io/gorm"
)

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
