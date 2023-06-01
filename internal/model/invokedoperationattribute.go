package model

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type InvokedOperationAttribute struct {
	gorm.Model
	Key   string `json:"key"`
	Value any    `json:"value" gorm:"json"`

	// FK to invokedoperation for 1..*
	InvokedOperationID uuid.UUID `json:"invokedOperationId" gorm:"type:uuid"`
}

func (a InvokedOperationAttribute) Set(key AttributeKey, v any) InvokedOperationAttribute {
	a.Key = string(key)
	a.Value = v
	return a
}

func NewAttribute(key AttributeKey, v any) InvokedOperationAttribute {
	return InvokedOperationAttribute{}.Set(key, v)
}
