package model

import "gorm.io/gorm"

type Deployment struct {
	gorm.Model
	Name     string `gorm:"unique"`
	Status   string
	Template map[string]any
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
