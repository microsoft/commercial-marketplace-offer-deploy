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
	Name     string `json:"name"`
	Status   string `json:"status"`
	Template any    `json:"template" gorm:"json"`
}
