package eventhook

import (
	"context"
	"fmt"
	"strings"

	log "github.com/sirupsen/logrus"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore/runtime"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/to"
	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/resources/armresources"
	"github.com/google/uuid"
	"github.com/microsoft/commercial-marketplace-offer-deploy/cmd/apiserver/azureevents"
	"github.com/microsoft/commercial-marketplace-offer-deploy/internal/data"
	"github.com/microsoft/commercial-marketplace-offer-deploy/internal/model"
	"github.com/microsoft/commercial-marketplace-offer-deploy/internal/structure"
	d "github.com/microsoft/commercial-marketplace-offer-deploy/pkg/deployment"
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
func (f *EventHookMessageFactory) convert(item *azureevents.ResourceEventSubject) (*sdk.EventHookMessage, error) {
	deployment, err := f.getRelatedDeployment(item)
	if err != nil {
		return nil, err
	}

	messageId, err := uuid.Parse(*item.Message.ID)
	if err != nil {
		return nil, err
	}

	eventData := azureevents.ResourceEventData{}
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
		Attempts:      int(invokedOperation.Attempts),
		Message:       *item.Message.Subject,
	}

	message := &sdk.EventHookMessage{
		Id:     messageId,
		Status: item.GetStatus(),
		Type:   f.getEventHookType(*item.Resource().Name, deployment),
		Data:   data,
	}

	// if this is a stage completed event, we need to set the stageId
	if message.Type == string(sdk.EventTypeStageCompleted) {
		for _, stage := range deployment.Stages {
			if *item.Resource().Name == stage.AzureDeploymentName {
				data.StageId = to.Ptr(stage.ID)
			}
		}
	}

	var stageId uuid.UUID
	lookupTags := item.LookupTags()

	if lookupTags[d.LookupTagKeyStageId] != nil {
		value := *lookupTags[d.LookupTagKeyStageId]
		stageId, _ = uuid.Parse(value)
	}
	message.SetSubject(deployment.ID, &stageId)

	return message, nil
}

// finds the deployment using the correlationId of the event which is ALWAYS 1-1 with the deployment.
//
//	remarks: the correlationId cannot be used to lookup the stage, since the stage will be a child deployment of the parent, and its correlationId still related to the parent
//			 that started the deployment
func (f *EventHookMessageFactory) getRelatedDeployment(item *azureevents.ResourceEventSubject) (*model.Deployment, error) {
	eventData := azureevents.ResourceEventData{}
	structure.Decode(item.Message.Data, &eventData)

	correlationId := eventData.CorrelationID
	resourceId := item.ResourceID()

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
					log.Debugf("correlationId [%s] matches [%s]", *item.Properties.CorrelationID, *item.Name)

					// try find the deployment using the name, then the stageId tag

					id, err := deployment.ParseAzureDeploymentName(*item.Name)

					// if we couldn't get the deployment by the resource deployment name, then it's a stage.
					if err != nil {
						log.Debugf("resource [%s] could not be found using correlationId [%s]", *item.Name, *item.Properties.CorrelationID)
					} else {
						return id, nil
					}

					if stageId, ok := item.Tags[string(d.LookupTagKeyId)]; ok {
						uuid, err := uuid.Parse(*stageId)
						if err == nil {
							deployment, _, err := f.findStage.Execute(uuid)

							if err == nil {
								return to.Ptr(int(deployment.ID)), nil
							}
						}
					}

					log.Errorf("error finding deploymentId for [%s] using correlationId [%s]", *item.Name, *item.Properties.CorrelationID)
				}
			}
		}
	}
	return nil, fmt.Errorf("deployment not found for correlationId: %s", correlationId)
}

func (f *EventHookMessageFactory) getEventHookType(resourceName string, deployment *model.Deployment) string {
	if resourceName == deployment.GetAzureDeploymentName() {
		return string(sdk.EventTypeDeploymentCompleted)
	} else {
		for _, stage := range deployment.Stages {
			if resourceName == stage.AzureDeploymentName {
				return string(sdk.EventTypeStageCompleted)
			}
		}
	}
	return string(sdk.EventTypeDeploymentEventReceived)
}

//endregion private methods
