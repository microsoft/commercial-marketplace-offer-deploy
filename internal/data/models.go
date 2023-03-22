package data

import (
	"gorm.io/gorm"
)

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
