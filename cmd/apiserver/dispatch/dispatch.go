package dispatch

import (
	"context"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	"github.com/google/uuid"
	"github.com/microsoft/commercial-marketplace-offer-deploy/internal/config"
	"github.com/microsoft/commercial-marketplace-offer-deploy/internal/data"
	"github.com/microsoft/commercial-marketplace-offer-deploy/internal/model"
	"github.com/microsoft/commercial-marketplace-offer-deploy/internal/operation"
	"github.com/microsoft/commercial-marketplace-offer-deploy/sdk"
	"gorm.io/gorm"
)

// Dispatch the invoke operation to the appropriate executor implemented in the operator
type OperatorDispatcher interface {
	Dispatch(ctx context.Context, command *DispatchInvokedOperation) (uuid.UUID, error)
}

// operator dispatcher
type dispatcher struct {
	db      *gorm.DB
	factory operation.Factory
}

func (h *dispatcher) Dispatch(ctx context.Context, command *DispatchInvokedOperation) (uuid.UUID, error) {
	err := validateOperationName(command.Name)
	if err != nil {
		return uuid.Nil, err
	}

	id, err := h.createOrUpdate(ctx, command)
	if err != nil {
		return uuid.Nil, err
	}

	invokedOperation, nil := h.factory.Create(ctx, id)
	if err != nil {
		return id, err
	}

	err = invokedOperation.Schedule()
	if err != nil {
		return invokedOperation.ID, err
	}

	return invokedOperation.ID, nil
}

func validateOperationName(name string) error {
	_, err := sdk.Type(name)
	return err
}

// createOrUpdate changes to the database
func (p *dispatcher) createOrUpdate(ctx context.Context, c *DispatchInvokedOperation) (uuid.UUID, error) {
	tx := p.db.Begin()

	invokedOperation := &model.InvokedOperation{}

	if c.OperationId != uuid.Nil {
		invokedOperation.ID = c.OperationId
		tx.First(invokedOperation)
	}

	invokedOperation.DeploymentId = uint(c.DeploymentId)
	invokedOperation.Name = c.Name
	invokedOperation.Status = string(sdk.StatusScheduled)
	invokedOperation.Parameters = c.Parameters
	invokedOperation.Attributes = c.Attributes
	invokedOperation.Retries = c.Retries

	tx.Save(invokedOperation)

	if tx.Error != nil {
		tx.Rollback()
		return uuid.Nil, tx.Error
	}

	tx.Commit()

	return invokedOperation.ID, nil
}

// region factory

func NewOperatorDispatcher(appConfig *config.AppConfig, credential azcore.TokenCredential) (OperatorDispatcher, error) {
	db := data.NewDatabase(appConfig.GetDatabaseOptions()).Instance()

	factory, err := operation.NewOperationFactory(appConfig, nil)
	if err != nil {
		return nil, err
	}

	return &dispatcher{
		db:      db,
		factory: factory,
	}, nil
}

//endregion factory
