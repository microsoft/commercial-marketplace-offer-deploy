package model

import (
	"github.com/google/uuid"
	"github.com/microsoft/commercial-marketplace-offer-deploy/sdk"
	"gorm.io/gorm"
)

type StageNotification struct {
	gorm.Model
	InvokedOperationId uuid.UUID                `json:"invokedOperationId" gorm:"type:uuid;not null"`
	CorrelationId      uuid.UUID                `json:"correlationId" gorm:"type:uuid;not null"`
	ResourceGroupName  string                   `json:"resourceGroupName" gorm:"not null"`
	Entries            []StageNotificationEntry `json:"entries" gorm:"json"`
	Done               bool                     `json:"done" gorm:"not null"`
}

type StageNotificationEntry struct {
	StageId uint                 `json:"stageId" gorm:"uuid;not null"`
	Message sdk.EventHookMessage `json:"message" gorm:"json;not null"`
	Sent    bool                 `json:"sent" gorm:"not null"`
}
