package operations

import (
	"context"

	"github.com/microsoft/commercial-marketplace-offer-deploy/internal/data"
	"github.com/microsoft/commercial-marketplace-offer-deploy/internal/hook"
	"gorm.io/gorm"
)

// deployment retry
// TODO: implement retry of an entire deployment and also a stage

type retryDeployment struct {
	db            *gorm.DB
	hookPublisher hook.Publisher
}

func (r *retryDeployment) Invoke(ctx context.Context, operation *data.InvokedOperation) error {
	return nil
}

//region factory

func NewRetryDeployment(db *gorm.DB, hookPublisher hook.Publisher) *retryDeployment {
	return &retryDeployment{
		db:            db,
		hookPublisher: hookPublisher,
	}
}

//endregion factory
