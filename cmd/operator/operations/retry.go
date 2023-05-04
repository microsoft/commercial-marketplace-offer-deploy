package operations

import (
	"context"

	"github.com/microsoft/commercial-marketplace-offer-deploy/internal/config"
	"github.com/microsoft/commercial-marketplace-offer-deploy/internal/data"
	log "github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

// deployment retry
// TODO: implement retry of an entire deployment and also a stage

type retryDeployment struct {
	db *gorm.DB
}

func (r *retryDeployment) Execute(ctx context.Context, operation *data.InvokedOperation) error {
	log.Infof("RetryDeploymentExecutor [%s]. enqueing to retry deployment", operation.ID)
	return nil
}

//region factory

func NewRetryDeploymentExecutor(appConfig *config.AppConfig) Executor {
	db := data.NewDatabase(appConfig.GetDatabaseOptions()).Instance()
	return &retryDeployment{
		db: db,
	}
}

//endregion factory
