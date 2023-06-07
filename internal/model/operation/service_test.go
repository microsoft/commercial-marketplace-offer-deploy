package operation

import (
	"context"
	"fmt"
	"testing"

	"github.com/google/uuid"
	"github.com/microsoft/commercial-marketplace-offer-deploy/internal/data"
	"github.com/microsoft/commercial-marketplace-offer-deploy/internal/model"
	"github.com/microsoft/commercial-marketplace-offer-deploy/sdk"
	log "github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func Test_Service_first(t *testing.T) {
	db := data.NewDatabase(&data.DatabaseOptions{UseInMemory: true}).Instance()
	correlationId := uuid.New()

	service := &OperationService{
		db:  db,
		log: log.WithContext(context.Background()),
		operation: &Operation{
			InvokedOperation: model.InvokedOperation{
				Status: sdk.StatusRunning.String(),
				Attributes: []model.InvokedOperationAttribute{
					model.NewAttribute(model.AttributeKeyCorrelationId, correlationId.String()),
				},
			},
		},
	}

	db.Save(&service.operation.InvokedOperation)

	service.id = service.operation.ID

	result, err := service.first()
	assert.NoError(t, err)
	assert.Len(t, result.Attributes, 1)
	assert.Equal(t, correlationId.String(), fmt.Sprintf("%s", result.Attributes[0].Value))
}

func Test_Service_notifyForStages_Only_Thats_Running(t *testing.T) {
	service := &OperationService{
		log: log.WithContext(context.Background()),
		operation: &Operation{
			InvokedOperation: model.InvokedOperation{
				Status: sdk.StatusRunning.String(),
			},
		},
	}

	err := service.notifyForStages()
	assert.NotNil(t, err)
	assert.NotEqual(t, err.Error(), "not a running deployment operation or not first attempt")
}

func Test_Service_notifyForStage(t *testing.T) {
	db := data.NewDatabase(&data.DatabaseOptions{UseInMemory: true}).Instance()
	db.Save(&model.Deployment{
		Stages: []model.Stage{
			{
				BaseWithGuidPrimaryKey: model.BaseWithGuidPrimaryKey{
					ID: uuid.New(),
				},
			},
		},
	})

	tx := db.Model(&model.Deployment{}).Preload("Stages").First(&model.Deployment{}, 1)
	require.Equal(t, int(tx.RowsAffected), int(1))

	correlationId := uuid.New()

	service := &OperationService{
		db:  db,
		log: log.WithContext(context.Background()),
		operation: &Operation{
			InvokedOperation: model.InvokedOperation{
				Status:       sdk.StatusRunning.String(),
				DeploymentId: 1,
				Retries:      1,
				Attempts:     1,
				Attributes: []model.InvokedOperationAttribute{
					model.NewAttribute(model.AttributeKeyCorrelationId, correlationId),
				},
			},
		},
		stageNotifier: &fakeStageNotifier{t: t},
	}

	require.Equal(t, "running", service.operation.Status)

	err := service.notifyForStages()
	assert.NoError(t, err)

	result := service.stageNotifier.(*fakeStageNotifier)

	assert.True(t, result.received)

	assert.Equal(t, service.operation.ID, result.notification.OperationId)
	assert.Equal(t, correlationId, result.notification.CorrelationId)
	assert.Len(t, result.notification.Entries, 1)
}

type fakeStageNotifier struct {
	t            *testing.T
	received     bool
	notification *model.StageNotification
}

func (f *fakeStageNotifier) Notify(ctx context.Context, notification *model.StageNotification) error {
	f.received = true
	f.notification = notification
	return nil
}
