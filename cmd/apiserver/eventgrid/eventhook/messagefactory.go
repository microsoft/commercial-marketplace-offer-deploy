package eventhook

import (
	"context"
	"fmt"
	"strings"

	"github.com/iancoleman/strcase"
	log "github.com/sirupsen/logrus"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore/arm"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/runtime"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/to"
	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/resources/armresources"
	"github.com/Azure/azure-sdk-for-go/services/eventgrid/2018-01-01/eventgrid"
	"github.com/google/uuid"
	eg "github.com/microsoft/commercial-marketplace-offer-deploy/cmd/apiserver/eventgrid"
	filtering "github.com/microsoft/commercial-marketplace-offer-deploy/cmd/apiserver/eventgrid/eventsfiltering"
	"github.com/microsoft/commercial-marketplace-offer-deploy/internal/model"
	"github.com/microsoft/commercial-marketplace-offer-deploy/internal/structure"
	d "github.com/microsoft/commercial-marketplace-offer-deploy/pkg/deployment"
	"github.com/microsoft/commercial-marketplace-offer-deploy/sdk"
	"gorm.io/gorm"
)

const (
	// we want to identify the write failure and success
	azureEventTypeResourceWriteFailure = "Microsoft.Resources.ResourceWriteFailure"
	azureEventTypeResourceWriteSuccess = "Microsoft.Resources.ResourceWriteSuccess"

	// identifies that we're dealing with an azure deployment resource, not just a resource being deployed as part of a deployment
	azureDeploymentResourceOperationName = "Microsoft.Resources/deployments/write"
)

// this factory is intented to create a list of WebHookEventMessages from a list of EventGridEventResources
// so the messages can be relayed via queue to be published to MODM consumer webhook subscription
type EventHookMessageFactory struct {
	client *armresources.DeploymentsClient
	filter filtering.EventGridEventFilter
	db     *gorm.DB
}

func NewEventHookMessageFactory(filter filtering.EventGridEventFilter, client *armresources.DeploymentsClient, db *gorm.DB) *EventHookMessageFactory {
	return &EventHookMessageFactory{
		client: client,
		filter: filter,
		db:     db,
	}
}

// Creates a list of messages from a list of EventGridEventResource
func (f *EventHookMessageFactory) Create(ctx context.Context, matchAny d.LookupTags, eventGridEvents []*eventgrid.Event) []*sdk.EventHookMessage {
	result := f.filter.Filter(ctx, matchAny, eventGridEvents)
	messages := []*sdk.EventHookMessage{}

	log.Debugf("factory received %d EventGridEvents, filtered to %d messages", len(eventGridEvents), len(result))

	for _, item := range result {
		message, err := f.convert(item)
		if err != nil {
			log.Errorf("failed to convert EventGridEventResource to WebHookEventMessage: %s", err.Error())
			continue
		}

		messages = append(messages, message)
	}

	return messages
}

//region private methods

// converts an EventGridEventResource to a WebHookEventMessage so it can be relayed via queue to be published to MODM consumer hook registration
func (f *EventHookMessageFactory) convert(item *eg.EventGridEventResource) (*sdk.EventHookMessage, error) {
	deployment, err := f.getRelatedDeployment(item)
	if err != nil {
		return nil, err
	}

	messageId, err := uuid.Parse(*item.Message.ID)
	if err != nil {
		return nil, err
	}

	eventData := eg.ResourceEventData{}
	structure.Decode(item.Message.Data, &eventData)

	// get related operation
	invokedOperation := &model.InvokedOperation{}
	f.db.Where("deployment_id = ? AND name = ?",
		deployment.ID,
		sdk.OperationDeploy,
	).First(invokedOperation)

	data := sdk.DeploymentEventData{
		DeploymentId:  int(deployment.ID),
		OperationId:   invokedOperation.ID,
		CorrelationId: to.Ptr(uuid.MustParse(eventData.CorrelationID)),
		Attempts:      invokedOperation.Attempts,
		Message:       *item.Message.Subject,
	}

	message := &sdk.EventHookMessage{
		Id:     messageId,
		Status: f.getStatus(*item.Message.EventType),
		Type:   f.getEventHookType(*item.Resource.Name, deployment),
		Data:   data,
	}

	// if this is a stage completed event, we need to set the stageId
	if message.Type == string(sdk.EventTypeStageCompleted) {
		for _, stage := range deployment.Stages {
			if *item.Resource.Name == stage.DeploymentName {
				data.StageId = to.Ptr(stage.ID)
			}
		}
	}

	var stageId uuid.UUID
	if item.Tags[d.LookupTagKeyStageId] != nil {
		value := *item.Tags[d.LookupTagKeyStageId]
		stageId, _ = uuid.Parse(value)
	}
	message.SetSubject(deployment.ID, &stageId)

	return message, nil
}

// finds the deployment using the correlationId of the event which is ALWAYS 1-1 with the deployment.
//
//	remarks: the correlationId cannot be used to lookup the stage, since the stage will be a child deployment of the parent, and its correlationId still related to the parent
//			 that started the deployment
func (f *EventHookMessageFactory) getRelatedDeployment(item *eg.EventGridEventResource) (*model.Deployment, error) {
	eventData := eg.ResourceEventData{}
	structure.Decode(item.Message.Data, &eventData)

	correlationId := eventData.CorrelationID
	resourceId, err := arm.ParseResourceID(*item.Resource.ID)
	if err != nil {
		return nil, err
	}

	pager := f.client.NewListByResourceGroupPager(resourceId.ResourceGroupName, nil)
	deploymentId, err := f.lookupDeploymentId(context.Background(), correlationId, pager)
	if err != nil {
		return nil, err
	}

	deployment := &model.Deployment{}
	tx := f.db.First(deployment, deploymentId)
	if tx.Error != nil {
		return nil, tx.Error
	}

	return deployment, nil
}

func (f *EventHookMessageFactory) lookupDeploymentId(ctx context.Context, correlationId string, pager *runtime.Pager[armresources.DeploymentsClientListByResourceGroupResponse]) (*int, error) {
	deployment := &model.Deployment{}

	for pager.More() {
		nextResult, err := pager.NextPage(ctx)
		if err != nil {
			return nil, err
		}

		if nextResult.DeploymentListResult.Value != nil {
			for _, item := range nextResult.DeploymentListResult.Value {
				correlationIdMatches := strings.EqualFold(*item.Properties.CorrelationID, correlationId)
				if correlationIdMatches {
					id, err := deployment.ParseAzureDeploymentName(*item.Name)
					if err != nil {
						fmt.Printf("correlation: resource/%s/correlationId/%s", *item.Name, *item.Properties.CorrelationID)

						// the name didn't match our pattern so we're not interested in this azure deployment, keep searching for a match
						// until we find 1=1 for our deployment (the top level "main deployment")
						continue
					} else {
						return id, nil
					}
				}
			}
		}
	}
	return nil, fmt.Errorf("deployment not found for correlationId: %s", correlationId)
}

func (f *EventHookMessageFactory) getStatus(eventType string) string {
	switch eventType {
	case azureEventTypeResourceWriteFailure:
		return sdk.StatusFailed.String()
	case azureEventTypeResourceWriteSuccess:
		return sdk.StatusSuccess.String()
	default:
		return strcase.ToCamel(eventType)
	}
}

func (f *EventHookMessageFactory) getEventHookType(resourceName string, deployment *model.Deployment) string {
	if resourceName == deployment.GetAzureDeploymentName() {
		return string(sdk.EventTypeDeploymentCompleted)
	} else {
		for _, stage := range deployment.Stages {
			if resourceName == stage.DeploymentName {
				return string(sdk.EventTypeStageCompleted)
			}
		}
	}
	return string(sdk.EventTypeDeploymentEventReceived)
}

//endregion private methods
