package model

import (
	"github.com/google/uuid"
	"github.com/microsoft/commercial-marketplace-offer-deploy/sdk"
	"gorm.io/gorm"
)

type StageNotification struct {
	gorm.Model
	OperationId       uuid.UUID                `json:"operationId" gorm:"type:uuid;not null"`
	CorrelationId     uuid.UUID                `json:"correlationId" gorm:"type:uuid;not null"`
	ResourceGroupName string                   `json:"resourceGroupName" gorm:"not null"`
	Entries           []StageNotificationEntry `json:"entries" gorm:"json"`
	Done              bool                     `json:"done" gorm:"not null"`
}

type StageNotificationEntry struct {
	StageId uuid.UUID            `json:"stageId" gorm:"uuid;not null"`
	Message sdk.EventHookMessage `json:"message" gorm:"json;not null"`
	Sent    bool                 `json:"sent" gorm:"not null"`
}
