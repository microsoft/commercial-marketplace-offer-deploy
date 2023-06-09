package notification

import (
	"context"
	"time"

	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/resources/armresources"
	"github.com/google/uuid"
	"github.com/microsoft/commercial-marketplace-offer-deploy/internal/hook"
	"github.com/microsoft/commercial-marketplace-offer-deploy/internal/model"
	"github.com/microsoft/commercial-marketplace-offer-deploy/pkg/deployment"
	log "github.com/sirupsen/logrus"
)

type StageNotificationHandlerFactoryFunc = NotificationHandlerFactoryFunc[model.StageNotification]

// package internal implementation of notification handler
type stageNotificationHandler struct {
	notify            hook.NotifyFunc
	deploymentsClient *armresources.DeploymentsClient
}

type resourceState struct {
	provisiningState armresources.ProvisioningState
	timestamp        *time.Time
}

type provisioningStates map[uuid.UUID]resourceState

func NewStageNotificationHandler(deploymentsClient *armresources.DeploymentsClient) NotificationHandler[model.StageNotification] {
	return &stageNotificationHandler{
		deploymentsClient: deploymentsClient,
		notify:            hook.Notify,
	}
}

func (h *stageNotificationHandler) Handle(context *NotificationHandlerContext[model.StageNotification]) {
	states, err := h.getStates(context.ctx, context.Notification)
	if err != nil {
		log.Errorf("failed to get states for notification [%s]", context.Notification.CorrelationId)
		context.Error(err)
		return
	}

	log.Tracef("total Stage Deployment states found: %d", len(states))
	if len(states) == 0 {
		context.Continue()
		return
	}

	entries := context.Notification.Entries
	for index, entry := range entries {
		if !entry.IsSent {
			state, ok := states[entry.StageId]
			if !ok {
				continue
			}
			if state.provisiningState != armresources.ProvisioningStateFailed && state.provisiningState != armresources.ProvisioningStateSucceeded {
				message := entry.Message
				data, err := message.StageEventData()
				if err != nil {
					log.Errorf("failed to get stage event data for stage [%s]", entry.StageId)
				}
				data.StartedAt = state.timestamp
				id, err := h.notify(context.ctx, &entry.Message)
				if err == nil {
					log.Tracef("notification sent for stage [%s] with id [%s]", entry.StageId, id)
					entry.Sent()
				}
			}
		}
		entries[index] = entry
	}

	if context.Notification.AllSent() {
		log.Tracef("all notifications sent for [%d]. Marking done", context.Notification.ID)

		context.Notification.Done()
		context.Done()
		return
	}

	context.Continue()
}

func (h *stageNotificationHandler) getStates(ctx context.Context, notification *model.StageNotification) (provisioningStates, error) {
	resources, err := h.getAzureDeploymentResources(ctx, notification)
	results := provisioningStates{}

	if err != nil || len(resources) == 0 {
		return results, err
	}

	for _, resource := range resources {
		value, ok := resource.Tags[string(deployment.LookupTagKeyId)]
		if !ok {
			continue
		}
		stageId, err := uuid.Parse(*value)
		if err != nil {
			log.Warnf("failed to uuid parse [%v] as modm.id on resource [%s]", *value, *resource.Name)
			continue
		}
		results[stageId] = resourceState{
			provisiningState: armresources.ProvisioningState(*resource.Properties.ProvisioningState),
			timestamp:        resource.Properties.Timestamp,
		}
	}
	return results, nil
}

// get by correlationId
func (h *stageNotificationHandler) getAzureDeploymentResources(ctx context.Context, notification *model.StageNotification) ([]*armresources.DeploymentExtended, error) {
	pager := h.deploymentsClient.NewListByResourceGroupPager(notification.ResourceGroupName, nil)

	resources := []*armresources.DeploymentExtended{}

	for pager.More() {
		nextResult, err := pager.NextPage(ctx)
		if err != nil {
			return nil, err
		}
		if nextResult.DeploymentListResult.Value != nil {
			for _, item := range nextResult.DeploymentListResult.Value {
				if item.Properties.CorrelationID == nil {
					continue
				}

				if uuid.MustParse(*item.Properties.CorrelationID) == notification.CorrelationId {
					resources = append(resources, item)
				}
			}
		}
	}
	return resources, nil
}
