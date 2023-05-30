package dispatch

import (
	"context"
	"strconv"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	"github.com/google/uuid"
	"github.com/microsoft/commercial-marketplace-offer-deploy/internal/config"
	"github.com/microsoft/commercial-marketplace-offer-deploy/internal/data"
	"github.com/microsoft/commercial-marketplace-offer-deploy/internal/hook"
	"github.com/microsoft/commercial-marketplace-offer-deploy/internal/messaging"
	"github.com/microsoft/commercial-marketplace-offer-deploy/internal/model"
	"github.com/microsoft/commercial-marketplace-offer-deploy/sdk"
	"gorm.io/gorm"
)

// Dispatch the invoke operation to the appropriate executor implemented in the operator
type OperatorDispatcher interface {
	Dispatch(ctx context.Context, command *DispatchInvokedOperation) (uuid.UUID, error)
}

// start deployment operation dispatcher
type dispatcher struct {
	db     *gorm.DB
	sender messaging.MessageSender
}

func (h *dispatcher) Dispatch(ctx context.Context, command *DispatchInvokedOperation) (uuid.UUID, error) {
	err := validateOperationName(*command.Request.Name)
	if err != nil {
		return uuid.Nil, err
	}

	invokedOperation, err := h.save(ctx, command)
	if err != nil {
		return uuid.Nil, err
	}

	err = h.send(ctx, invokedOperation.ID)

	if err != nil {
		return uuid.Nil, err
	}

	h.addEventHook(ctx, invokedOperation)

	return invokedOperation.ID, nil
}

func validateOperationName(name string) error {
	_, err := sdk.Type(name)
	return err
}

// save changes to the database
func (p *dispatcher) save(ctx context.Context, c *DispatchInvokedOperation) (*model.InvokedOperation, error) {
	tx := p.db.Begin()

	invokedOperation := &model.InvokedOperation{
		DeploymentId: uint(c.DeploymentId),
		Name:         *c.Request.Name,
		Status:       string(sdk.StatusScheduled),
		Parameters:   c.Request.Parameters.(map[string]interface{}),
	}

	invokedOperation.Retries = int(*c.Request.Retries)
	if *c.Request.Retries <= 0 {
		invokedOperation.Retries = 1
	}

	tx.Save(invokedOperation)

	if tx.Error != nil {
		tx.Rollback()
		return nil, tx.Error
	}

	tx.Commit()

	return invokedOperation, nil
}

// send the message on the queue
func (h *dispatcher) send(ctx context.Context, operationId uuid.UUID) error {
	message := messaging.ExecuteInvokedOperation{OperationId: operationId}

	results, err := h.sender.Send(ctx, string(messaging.QueueNameOperations), message)
	if err != nil {
		return err
	}

	if len(results) == 1 && results[0].Error != nil {
		return results[0].Error
	}
	return nil
}

func (h *dispatcher) addEventHook(ctx context.Context, invokedOperation *model.InvokedOperation) error {
	return hook.Add(ctx, &sdk.EventHookMessage{
		Status:  invokedOperation.Status,
		Type:    string(sdk.EventTypeDeploymentOperationReceived),
		Subject: "/deployments/" + strconv.Itoa(int(invokedOperation.DeploymentId)),
		Data: &sdk.DeploymentEventData{
			DeploymentId: int(invokedOperation.DeploymentId),
			OperationId:  invokedOperation.ID,
			Message:      "",
		},
	})
}

// region factory

func NewOperatorDispatcher(appConfig *config.AppConfig, credential azcore.TokenCredential) (OperatorDispatcher, error) {
	db := data.NewDatabase(appConfig.GetDatabaseOptions()).Instance()
	sender, err := newMessageSender(appConfig, credential)
	if err != nil {
		return nil, err
	}

	return &dispatcher{
		db:     db,
		sender: sender,
	}, nil
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
