package model

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type AttributeKey string

const (
	AttributeKeyCorrelationId AttributeKey = "correlationId"
	AttributeKeyResumeToken   AttributeKey = "resumeToken"
)

type InvokedOperationAttribute struct {
	gorm.Model
	Key                string    `json:"key" gorm:"uniqueIndex:composite_index;index;not null"`
	Value              any       `json:"value" gorm:"json"`
	InvokedOperationID uuid.UUID `json:"invokedOperationId" gorm:"type:uuid;uniqueIndex:composite_index;index;not null"`
}

func (a InvokedOperationAttribute) Set(key AttributeKey, v any) InvokedOperationAttribute {
	a.Key = string(key)
	a.Value = v
	return a
}

func NewAttribute(key AttributeKey, v any) InvokedOperationAttribute {
	return InvokedOperationAttribute{}.Set(key, v)
}
