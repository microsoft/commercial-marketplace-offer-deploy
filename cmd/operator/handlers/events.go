package handlers

import (
	"context"
	"encoding/json"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	"github.com/google/uuid"
	"github.com/microsoft/commercial-marketplace-offer-deploy/internal/config"
	"github.com/microsoft/commercial-marketplace-offer-deploy/internal/data"
	"github.com/microsoft/commercial-marketplace-offer-deploy/internal/hook"
	"github.com/microsoft/commercial-marketplace-offer-deploy/internal/messaging"
	"github.com/microsoft/commercial-marketplace-offer-deploy/internal/structure"
	"github.com/microsoft/commercial-marketplace-offer-deploy/sdk"
	log "github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

type eventsMessageHandler struct {
	db        *gorm.DB
	publisher hook.Publisher
	sender    messaging.MessageSender
}

func (h *eventsMessageHandler) Handle(message *sdk.EventHookMessage, context messaging.MessageHandlerContext) error {
	bytes, _ := json.Marshal(message)
	log.WithField("eventHookMessage", string(bytes)).Debug("Events Handler excuting")

	if h.shouldRetryIfDeployment(message) {
		h.retryDeployment(context.Context(), message)
	}

	err := h.publisher.Publish(message)
	return err
}

func (h *eventsMessageHandler) shouldRetryIfDeployment(message *sdk.EventHookMessage) bool {
	failedStatus := message != nil && message.Status == sdk.StatusFailed.String()
	deploymentTypes := message.Type == string(sdk.EventTypeDeploymentCompleted) || message.Type == string(sdk.EventTypeDeploymentRetried)

	return failedStatus && deploymentTypes
}

func (h *eventsMessageHandler) retryDeployment(ctx context.Context, message *sdk.EventHookMessage) error {
	log.Infof("EventHookMessage [%s]. enqueing to retry deployment", message.Id)

	invokedOperation, err := h.update(message)
	if err != nil {
		return err
	}

	if invokedOperation.IsRetriable() {
		results, err := h.sender.Send(ctx, string(messaging.QueueNameOperations), messaging.ExecuteInvokedOperation{OperationId: invokedOperation.ID})
		if err != nil {
			return err
		}

		if len(results) == 1 && results[0].Error != nil {
			return results[0].Error
		}

		err = hook.Add(ctx, &sdk.EventHookMessage{
			Status: invokedOperation.Status,
			Type:   string(sdk.EventTypeDeploymentCompleted),
			Data: sdk.DeploymentEventData{
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
	}
	return nil
}

// updates the invoked operation
func (h *eventsMessageHandler) update(message *sdk.EventHookMessage) (*data.InvokedOperation, error) {
	eventData := &sdk.DeploymentEventData{}
	structure.Decode(message.Data, &eventData)

	invokedOperation := &data.InvokedOperation{}
	h.db.First(&invokedOperation, eventData.OperationId)

	//update the status to scheduled
	if invokedOperation.IsRetriable() {
		invokedOperation.Status = string(sdk.StatusScheduled)
	} else {
		invokedOperation.Status = message.Status
	}

	if eventData.CorrelationId != nil && *eventData.CorrelationId != uuid.Nil {
		invokedOperation.CorrelationId = eventData.CorrelationId
	}

	h.db.Save(invokedOperation)

	return invokedOperation, h.db.Error
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
