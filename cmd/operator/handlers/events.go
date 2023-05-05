package handlers

import (
	"context"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	"github.com/microsoft/commercial-marketplace-offer-deploy/internal/config"
	"github.com/microsoft/commercial-marketplace-offer-deploy/internal/data"
	"github.com/microsoft/commercial-marketplace-offer-deploy/internal/hook"
	"github.com/microsoft/commercial-marketplace-offer-deploy/internal/messaging"
	"github.com/microsoft/commercial-marketplace-offer-deploy/pkg/events"
	"github.com/microsoft/commercial-marketplace-offer-deploy/pkg/operation"
	"github.com/mitchellh/mapstructure"
	log "github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

type eventsMessageHandler struct {
	db        *gorm.DB
	publisher hook.Publisher
	sender    messaging.MessageSender
}

func (h *eventsMessageHandler) Handle(message *events.EventHookMessage, context messaging.MessageHandlerContext) error {
	log.WithField("eventHookMessage", message).Debug("Received event hook message")

	if h.shouldRetryIfDeployment(message) {
		h.retryDeployment(context.Context(), message)
	}
	err := h.publisher.Publish(message)
	return err
}

func (h *eventsMessageHandler) shouldRetryIfDeployment(message *events.EventHookMessage) bool {
	failedStatus := message != nil && message.Status == operation.StatusFailed.String()
	deploymentTypes := message.Type == string(events.EventTypeDeploymentCompleted) || message.Type == string(events.EventTypeDeploymentRetried)

	return failedStatus && deploymentTypes
}

func (h *eventsMessageHandler) retryDeployment(ctx context.Context, message *events.EventHookMessage) error {
	log.Infof("EventHookMessage [%s]. enqueing to retry deployment", message.Id)

	eventData := events.DeploymentEventData{}
	mapstructure.Decode(message, &eventData)

	invokedOperation := &data.InvokedOperation{}
	h.db.First(&invokedOperation, eventData.OperationId)

	//update the status to scheduled
	invokedOperation.Status = string(operation.StatusScheduled)
	h.db.Save(&invokedOperation)

	results, err := h.sender.Send(ctx, string(messaging.QueueNameOperations), messaging.ExecuteInvokedOperation{OperationId: invokedOperation.ID})
	if err != nil {
		return err
	}

	if len(results) == 1 && results[0].Error != nil {
		return results[0].Error
	}

	err = hook.Add(ctx, &events.EventHookMessage{
		Status: invokedOperation.Status,
		Type:   string(events.EventTypeDeploymentCompleted),
		Data: events.DeploymentEventData{
			DeploymentId: int(invokedOperation.DeploymentId),
			OperationId:  invokedOperation.ID,
			Message:      "Deployment is being retried",
			Attempts:     invokedOperation.Attempts,
		},
	})

	if err != nil {
		log.Errorf("Failed to add event hook message. Error: %v", err)
		return err
	}
	return nil
}

//region factory

func NewEventsMessageHandler(appConfig *config.AppConfig, credential azcore.TokenCredential) (*eventsMessageHandler, error) {
	db := data.NewDatabase(appConfig.GetDatabaseOptions()).Instance()
	publisher := newWebHookPublisher(db)
	sender, err := newMessageSender(appConfig, credential)
	if err != nil {
		return nil, err
	}
	return &eventsMessageHandler{
		db:        db,
		publisher: publisher,
		sender:    sender,
	}, nil
}

func newWebHookPublisher(db *gorm.DB) hook.Publisher {
	subscriptionsProvider := hook.NewEventHooksProvider(db)
	publisher := hook.NewEventHookPublisher(subscriptionsProvider)
	return publisher
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
