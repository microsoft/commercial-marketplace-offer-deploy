package operations

import (
	"context"
	"log"

	//"time"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	"github.com/google/uuid"
	"github.com/microsoft/commercial-marketplace-offer-deploy/internal/config"
	"github.com/microsoft/commercial-marketplace-offer-deploy/internal/data"
	"github.com/microsoft/commercial-marketplace-offer-deploy/internal/messaging"
	"github.com/microsoft/commercial-marketplace-offer-deploy/internal/utils"
	"github.com/microsoft/commercial-marketplace-offer-deploy/pkg/events"
	"github.com/microsoft/commercial-marketplace-offer-deploy/pkg/operations"

	"gorm.io/gorm"
)

type InvokeOperationProcessor interface {
	Process(ctx context.Context, command *InvokeOperationCommand) (uuid.UUID, error)
}

// start deployment operation processor
type processor struct {
	db     *gorm.DB
	sender messaging.MessageSender
}

func (h *processor) Process(ctx context.Context, command *InvokeOperationCommand) (uuid.UUID, error) {
	err := validateOperationName(*command.Request.Name)
	if err != nil {
		return uuid.Nil, err
	}

	id, err := h.save(command)
	if err != nil {
		return uuid.Nil, err
	}

	err = h.send(ctx, id)

	if err != nil {
		return uuid.Nil, err
	}

	return id, nil
}

func validateOperationName(name string) error {
	_, err := operations.Type(name)
	return err
}

// save changes to the database
func (h *processor) save(c *InvokeOperationCommand) (uuid.UUID, error) {
	tx := h.db.Begin()

	deployment := &data.Deployment{}
	tx.First(&deployment, c.DeploymentId)
	deployment.Status = events.DeploymentPendingEventType.String()
	tx.Save(deployment)

	invokedOperation := &data.InvokedOperation{
		DeploymentId: uint(c.DeploymentId),
		Name:         *c.Request.Name,
		Parameters:   c.Request.Parameters.(map[string]interface{}),
	}
	tx.Save(invokedOperation)

	if tx.Error != nil {
		tx.Rollback()
		return uuid.Nil, tx.Error
	}

	tx.Commit()

	return invokedOperation.ID, nil
}

// send the message on the queue
func (h *processor) send(ctx context.Context, operationId uuid.UUID) error {
	message := messaging.InvokedOperationMessage{OperationId: operationId.String()}
	log.Printf("sending message from api/operations/processor.go - message: %v", message)
	results, err := h.sender.Send(ctx, string(messaging.QueueNameOperations), message)
	if err != nil {
		return err
	}

	if len(results) == 1 && results[0].Error != nil {
		return results[0].Error
	}
	return nil
}

// region factory

func NewInvokeOperationProcessor(appConfig *config.AppConfig, credential azcore.TokenCredential) (InvokeOperationProcessor, error) {
	errors := []string{}

	db := data.NewDatabase(appConfig.GetDatabaseOptions()).Instance()
	sender, err := newMessageSender(appConfig, credential)
	if err != nil {
		errors = append(errors, err.Error())
	}

	if len(errors) > 0 {
		return nil, utils.NewAggregateError(errors)
	}

	processor := &processor{
		db:     db,
		sender: sender,
	}
	return processor, nil
}

func newMessageSender(appConfig *config.AppConfig, credential azcore.TokenCredential) (messaging.MessageSender, error) {
	sender, err := messaging.NewServiceBusMessageSender(credential, messaging.MessageSenderOptions{
		SubscriptionId:          appConfig.Azure.SubscriptionId,
		Location:                appConfig.Azure.Location,
		ResourceGroupName:       appConfig.Azure.ResourceGroupName,
		FullyQualifiedNamespace: appConfig.Azure.GetFullQualifiedNamespace(),
	})
	if err != nil {
		return nil, err
	}

	return sender, nil
}

//endregion factory
