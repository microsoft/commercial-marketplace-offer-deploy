package app

import (
	"context"

	"github.com/google/uuid"
	"github.com/microsoft/commercial-marketplace-offer-deploy/internal/config"
	"github.com/microsoft/commercial-marketplace-offer-deploy/internal/data"
	"github.com/microsoft/commercial-marketplace-offer-deploy/internal/model"
	"github.com/microsoft/commercial-marketplace-offer-deploy/internal/model/operation"
	"github.com/microsoft/commercial-marketplace-offer-deploy/internal/tasks"
	"github.com/microsoft/commercial-marketplace-offer-deploy/internal/utils"
	"github.com/microsoft/commercial-marketplace-offer-deploy/sdk"
	"gorm.io/gorm"
)

// creates new task that continues operations if the operator was terminated
func newContinueOperationsTask(appConfig *config.AppConfig) tasks.Task {
	action := func(ctx context.Context) error {
		db := data.NewDatabase(appConfig.GetDatabaseOptions()).Instance()
		ids := getRunningOperationIds(db)

		if len(ids) == 0 {
			return nil
		}

		scheduler, err := operation.NewSchedulerFromConfig(appConfig)
		if err != nil {
			return err
		}

		errorMessages := []string{}

		for _, id := range ids {
			err := scheduler.Schedule(ctx, id)
			if err != nil {
				errorMessages = append(errorMessages, err.Error())
			}
		}
		return utils.NewAggregateError(errorMessages)
	}
	return tasks.NewTask("Continue Operations Task", action)

}

func getRunningOperationIds(db *gorm.DB) []uuid.UUID {
	ids := []uuid.UUID{}
	db.Model(&model.InvokedOperation{}).Where("status = ?", sdk.StatusRunning).Select("operation_id").Find(&ids)
	return ids
}
