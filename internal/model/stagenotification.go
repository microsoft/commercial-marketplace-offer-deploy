package model

import (
	"time"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore/to"
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
	IsDone            bool                     `json:"done" gorm:"not null"`
}

type StageNotificationEntry struct {
	StageId uuid.UUID            `json:"stageId" gorm:"uuid;not null"`
	Message sdk.EventHookMessage `json:"message" gorm:"json;not null"`
	IsSent  bool                 `json:"sent" gorm:"not null"`
	SentAt  *time.Time           `json:"sentAt" gorm:"null"`
}

func (n *StageNotificationEntry) Sent() {
	n.IsSent = true
	n.SentAt = to.Ptr(time.Now().UTC())
}

func (n *StageNotification) Done() {
	n.IsDone = true
}

// whether all entries have been sent
func (n *StageNotification) AllSent() bool {
	for _, entry := range n.Entries {
		if !entry.IsSent {
			return false
		}
	}
	return true
}
