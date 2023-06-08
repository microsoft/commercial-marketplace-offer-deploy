package eventhook

import (
	"context"

	log "github.com/sirupsen/logrus"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore/to"
	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/resources/armresources"
	"github.com/google/uuid"
	"github.com/microsoft/commercial-marketplace-offer-deploy/cmd/apiserver/azureevents"
	"github.com/microsoft/commercial-marketplace-offer-deploy/internal/data"
	"github.com/microsoft/commercial-marketplace-offer-deploy/sdk"
	"gorm.io/gorm"
)

// this factory is intented to create a list of WebHookEventMessages from a list of ResourceEventSubjects
// so the messages can be relayed via queue to be published to MODM consumer webhook subscription
type EventHookMessageFactory struct {
	client    *armresources.DeploymentsClient
	db        *gorm.DB
	findStage *data.StageQuery
}

func NewEventHookMessageFactory(client *armresources.DeploymentsClient, db *gorm.DB) *EventHookMessageFactory {
	return &EventHookMessageFactory{
		client:    client,
		db:        db,
		findStage: data.NewStageQuery(db),
	}
}

// Creates a list of messages from a list of ResourceEventSubject
func (f *EventHookMessageFactory) Create(ctx context.Context, resources []*azureevents.ResourceEventSubject) []*sdk.EventHookMessage {
	messages := []*sdk.EventHookMessage{}
	for _, item := range resources {
		message, err := f.convert(item)
		if err != nil {
			log.Errorf("failed to convert ResourceEventSubject to WebHookEventMessage: %s", err.Error())
			continue
		}

		messages = append(messages, message)
	}

	return messages
}

//region private methods

// converts an ResourceEventSubject to a WebHookEventMessage so it can be relayed via queue to be published to MODM consumer hook registration
func (f *EventHookMessageFactory) convert(subject *azureevents.ResourceEventSubject) (*sdk.EventHookMessage, error) {
	operation := subject.Operation()
	deployment := operation.Deployment()

	messageId, err := uuid.Parse(*subject.Message.ID)
	if err != nil {
		return nil, err
	}

	stageId, _ := subject.StageId()

	data := sdk.DeploymentEventData{
		EventData: sdk.EventData{
			DeploymentId: int(deployment.ID),
			OperationId:  operation.ID,
			Attempts:     int(operation.Attempts),
		},
		StageId:       stageId,
		CorrelationId: to.Ptr(uuid.MustParse(subject.CorrelationID())),
	}

	firstAttempt := operation.LatestResult()
	if firstAttempt != nil {
		data.CompletedAt = firstAttempt.CompletedAt.UTC()
	}

	latestAttempt := operation.LatestResult()
	if latestAttempt != nil {
		data.CompletedAt = latestAttempt.CompletedAt.UTC()
	}

	if operation.IsCompleted() {
		completedAt, _ := operation.CompletedAt()
		data.CompletedAt = *completedAt
	}

	message := &sdk.EventHookMessage{
		Id:     messageId,
		Status: subject.GetStatus(),
		Type:   subject.GetType(),
		Data:   data,
	}
	message.SetSubject(deployment.ID, stageId)

	return message, nil
}

//endregion private methods
