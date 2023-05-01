package operations

import (
	"context"

	"github.com/microsoft/commercial-marketplace-offer-deploy/internal/data"
	"github.com/microsoft/commercial-marketplace-offer-deploy/internal/hook"
	"gorm.io/gorm"
)

type retryStage struct {
	db            *gorm.DB
	hookPublisher hook.Publisher
}

func (r *retryStage) Invoke(ctx context.Context, operation *data.InvokedOperation) error {
	return nil
}

//region factory

func NewRetryStage(db *gorm.DB, hookPublisher hook.Publisher) *retryStage {
	return &retryStage{
		db:            db,
		hookPublisher: hookPublisher,
	}
}

//endregion factory
