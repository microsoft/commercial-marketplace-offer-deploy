package operation

import (
	"testing"

	"github.com/google/uuid"
	"github.com/microsoft/commercial-marketplace-offer-deploy/internal/data"
	"github.com/microsoft/commercial-marketplace-offer-deploy/internal/model"
	"github.com/microsoft/commercial-marketplace-offer-deploy/sdk"
	"github.com/stretchr/testify/assert"
	"gorm.io/gorm/clause"
)

func Test_Operation_CorrelationId(t *testing.T) {
	db := data.NewDatabase(&data.DatabaseOptions{UseInMemory: true}).Instance()
	correlationId := uuid.New()

	operation := &model.InvokedOperation{
		Status: sdk.StatusRunning.String(),
		Attributes: []model.InvokedOperationAttribute{
			model.NewAttribute(model.AttributeKeyCorrelationId, correlationId),
		},
	}

	db.Save(operation)

	result := &model.InvokedOperation{}
	db.Preload(clause.Associations).First(result, operation.ID)

	assert.Len(t, result.Attributes, 1)

	resultId, err := result.CorrelationId()
	assert.NoError(t, err)

	assert.Equal(t, correlationId, *resultId)
}

func Test_Operation_CorrelationId_Pointer(t *testing.T) {
	db := data.NewDatabase(&data.DatabaseOptions{UseInMemory: true}).Instance()

	var correlationId *uuid.UUID
	id := uuid.New()
	correlationId = &id

	operation := &model.InvokedOperation{
		Status: sdk.StatusRunning.String(),
		Attributes: []model.InvokedOperationAttribute{
			model.NewAttribute(model.AttributeKeyCorrelationId, correlationId),
		},
	}

	db.Save(operation)

	result := &model.InvokedOperation{}
	db.Preload(clause.Associations).First(result, operation.ID)

	assert.Len(t, result.Attributes, 1)

	resultId, _ := result.CorrelationId()
	assert.NotEqual(t, correlationId, *resultId)
}

func Test_Operation_CorrelationId_Attrs_Null(t *testing.T) {
	db := data.NewDatabase(&data.DatabaseOptions{UseInMemory: true}).Instance()

	correlationId := uuid.New()

	operation := &model.InvokedOperation{
		Status: sdk.StatusRunning.String(),
		Attributes: []model.InvokedOperationAttribute{
			model.NewAttribute(model.AttributeKeyCorrelationId, correlationId),
		},
	}

	db.Save(operation)

	result := &model.InvokedOperation{}
	db.Preload(clause.Associations).First(result, operation.ID)
	result.Attributes = nil

	resultId, _ := result.CorrelationId()
	assert.NotEqual(t, correlationId, resultId)
}
