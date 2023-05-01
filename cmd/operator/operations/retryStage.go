package operations

import (
	"context"

	"github.com/microsoft/commercial-marketplace-offer-deploy/internal/config"
	"github.com/microsoft/commercial-marketplace-offer-deploy/internal/data"
	"gorm.io/gorm"
)

type retryStage struct {
	db *gorm.DB
}

func (r *retryStage) Execute(ctx context.Context, operation *data.InvokedOperation) error {
	return nil
}

//region factory

func NewRetryStageExecutor(appConfig *config.AppConfig) Executor {
	db := data.NewDatabase(appConfig.GetDatabaseOptions()).Instance()
	return &retryStage{
		db: db,
	}
}

//endregion factory
