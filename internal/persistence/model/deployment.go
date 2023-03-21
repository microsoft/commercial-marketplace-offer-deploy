package model

import "gorm.io/gorm"

type Deployment struct {
	gorm.Model
	Name     string `gorm:"unique"`
	Status   string
	Template *DeploymentTemplate
	Stages   []Stage `gorm:"embedded"`
}

type Stage struct {
	gorm.Model
	Name   string `gorm:"unique"`
	Status string
}

type DeploymentTemplate struct {
	FilePath string
}
