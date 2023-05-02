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
	"github.com/microsoft/commercial-marketplace-offer-deploy/pkg/events"
	"github.com/microsoft/commercial-marketplace-offer-deploy/pkg/operations"
	log "github.com/sirupsen/logrus"
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

	id, err := h.save(ctx, command)
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
func (p *dispatcher) save(ctx context.Context, c *DispatchInvokedOperation) (uuid.UUID, error) {
	tx := p.db.Begin()

	deployment := &data.Deployment{}
	tx.First(&deployment, c.DeploymentId)

	// TODO: update deployment status depending on what the operation is
	deployment.Status = string(events.StatusScheduled)
	tx.Save(deployment)

	invokedOperation := &data.InvokedOperation{
		DeploymentId: uint(c.DeploymentId),
		Name:         *c.Request.Name,
		Status:       events.StatusAccepted.String(),
		Parameters:   c.Request.Parameters.(map[string]interface{}),
	}

	invokedOperation.Retries = int(*c.Request.Retries)
	log.Debugf("Retries is received to %d", *c.Request.Retries)
	if *c.Request.Retries <= 0 {
		log.Debug("Retries is not set, defaulting to 1")
		invokedOperation.Retries = 1
	}
	log.Debugf("Retries is set to %d", invokedOperation.Retries)
	tx.Save(invokedOperation)

	if tx.Error != nil {
		tx.Rollback()
		return uuid.Nil, tx.Error
	}

	tx.Commit()
	p.addEventHook(ctx, c.DeploymentId, invokedOperation.Name, invokedOperation.Status)

	return invokedOperation.ID, nil
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

func (h *dispatcher) addEventHook(ctx context.Context, deploymentId int, status string, operationName string) error {
	return hook.Add(&events.EventHookMessage{
		Id:      uuid.New(),
		Status:  status,
		Type:    string(events.EventTypeDeploymentOperationReceived),
		Subject: "/deployments/" + strconv.Itoa(deploymentId) + "/operations/" + operationName,
		Data: &events.DeploymentEventData{
			DeploymentId: deploymentId,
			Message:      "Operation " + operationName + " accepted",
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
