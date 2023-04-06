package operations

import (
	"context"
	"time"

	//"time"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/to"
	"github.com/google/uuid"
	"github.com/microsoft/commercial-marketplace-offer-deploy/internal/config"
	"github.com/microsoft/commercial-marketplace-offer-deploy/internal/data"
	"github.com/microsoft/commercial-marketplace-offer-deploy/internal/messaging"
	"github.com/microsoft/commercial-marketplace-offer-deploy/internal/utils"
	"github.com/microsoft/commercial-marketplace-offer-deploy/pkg/api"
	"github.com/microsoft/commercial-marketplace-offer-deploy/pkg/events"
	"github.com/microsoft/commercial-marketplace-offer-deploy/pkg/operations"

	"gorm.io/gorm"
)

// start deployment operation handler
type startDeploymentOperation struct {
	db     *gorm.DB
	sender messaging.MessageSender
}

func (h *startDeploymentOperation) Handle(c *DeploymentOperationContext) (*api.InvokedDeploymentOperation, error) {
	id, err := h.save(c)
	if err != nil {
		return nil, err
	}

	ctx := c.HttpContext.Request().Context()
	err = h.send(ctx, id)

	if err != nil {
		return nil, err
	}

	return &api.InvokedDeploymentOperation{
		ID:         to.Ptr(id.String()),
		InvokedOn:  to.Ptr(time.Now().UTC()),
		Name:       to.Ptr(c.Name.String()),
		Parameters: c.Parameters,
		Result:     operations.OperationResultAccepted,
		Status:     to.Ptr(events.DeploymentPendingEventType.String()),
	}, nil
}

// save changes to the database
func (h *startDeploymentOperation) save(c *DeploymentOperationContext) (uuid.UUID, error) {
	tx := h.db.Begin()

	deployment := &data.Deployment{}
	tx.First(&deployment, c.DeploymentId)
	deployment.Status = events.DeploymentPendingEventType.String()
	tx.Save(deployment)

	invokedOperation := &data.InvokedOperation{
		DeploymentId: uint(c.DeploymentId),
		Name:         c.Name.String(),
		Parameters:   c.Parameters,
	}
	tx.Save(&invokedOperation)

	if tx.Error != nil {
		tx.Rollback()
		return uuid.Nil, tx.Error
	}

	tx.Commit()

	return invokedOperation.ID, nil
}

// send the message on the queue
func (h *startDeploymentOperation) send(ctx context.Context, operationId uuid.UUID) error {
	message := messaging.InvokedOperationMessage{OperationId: operationId.String()}

	results, err := h.sender.Send(ctx, messaging.OperatorQueue, message)
	if err != nil {
		return err
	}

	if len(results) == 1 && results[0].Error != nil {
		return results[0].Error
	}
	return nil
}

// region factory

func NewStartDeploymentOperationHandler(appConfig *config.AppConfig, credential azcore.TokenCredential) DeploymentOperationHandlerFunc {
	return func(c *DeploymentOperationContext) (*api.InvokedDeploymentOperation, error) {
		errors := []string{}

		db := data.NewDatabase(appConfig.GetDatabaseOptions()).Instance()
		sender, err := newMessageSender(appConfig, credential)
		if err != nil {
			errors = append(errors, err.Error())
		}

		if len(errors) > 0 {
			return nil, utils.NewAggregateError(errors)
		}

		handler := &startDeploymentOperation{
			db:     db,
			sender: sender,
		}
		return handler.Handle(c)
	}
}

func newMessageSender(appConfig *config.AppConfig, credential azcore.TokenCredential) (messaging.MessageSender, error) {
	sender, err := messaging.NewServiceBusMessageSender(messaging.MessageSenderOptions{
		SubscriptionId:      appConfig.Azure.SubscriptionId,
		Location:            appConfig.Azure.Location,
		ResourceGroupName:   appConfig.Azure.ResourceGroupName,
		ServiceBusNamespace: appConfig.Azure.ServiceBusNamespace,
	}, credential)
	if err != nil {
		return nil, err
	}

	return sender, nil
}

//endregion factory
