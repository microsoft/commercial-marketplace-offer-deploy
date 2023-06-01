package data

import (
	"github.com/google/uuid"
	"github.com/microsoft/commercial-marketplace-offer-deploy/internal/model"
	"github.com/microsoft/commercial-marketplace-offer-deploy/sdk"
	"gorm.io/gorm"
)

type InvokedOperationQuery struct {
	db *gorm.DB
}

// finds the first invoked operation for a retry stage by id and correlation id
// if found returns the operation and nil error, otherwise returns nil, error
func (q *InvokedOperationQuery) First(stageId uuid.UUID, correlationId string) (*model.InvokedOperation, error) {
	var operations []model.InvokedOperation
	result := q.db.Model(&model.InvokedOperation{}).Where("name = ?", string(sdk.OperationRetryStage)).Find(&operations)

	if result.Error != nil {
		return nil, result.Error
	}

	for _, operation := range operations {
		if operation.Parameters == nil {
			continue
		}

		if value, ok := operation.Parameters["stageId"]; ok {
			if value == stageId.String() {
				if operation.Attributes == nil {
					continue
				}
				if value, ok := operation.AttributeValue(model.AttributeKeyCorrelationId); ok {
					if value == correlationId {
						return &operation, nil
					}
				}
			}
		}
	}
	return nil, nil
}

func NewInvokedOperationQuery(db *gorm.DB) *InvokedOperationQuery {
	return &InvokedOperationQuery{
		db: db,
	}
}
