package notification

import (
	"context"

	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/resources/armresources"
	"github.com/google/uuid"
	"github.com/microsoft/commercial-marketplace-offer-deploy/internal/hook"
	"github.com/microsoft/commercial-marketplace-offer-deploy/internal/model"
	"github.com/microsoft/commercial-marketplace-offer-deploy/pkg/deployment"
	log "github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

type StageNotificationHandlerFactoryFunc = NotificationHandlerFactoryFunc[model.StageNotification]

// package internal implementation of notification handler
type stageNotificationHandler struct {
	db                *gorm.DB
	notify            hook.NotifyFunc
	deploymentsClient *armresources.DeploymentsClient
}

type provisioningStates map[uuid.UUID]armresources.ProvisioningState

func NewStageNotificationHandler(db *gorm.DB, deploymentsClient *armresources.DeploymentsClient) NotificationHandler[model.StageNotification] {
	return &stageNotificationHandler{
		db:                db,
		deploymentsClient: deploymentsClient,
		notify:            hook.Notify,
	}
}

func (h *stageNotificationHandler) Handle(context *NotificationHandlerContext[model.StageNotification]) {
	states, err := h.getStates(context.ctx, context.Notification)
	if err != nil {
		context.Error(err)
		return
	}

	if len(states) == 0 {
		context.Continue()
		return
	}

	entries := context.Notification.Entries
	for _, entry := range entries {
		if !entry.IsSent {
			state, ok := states[entry.StageId]
			if !ok {
				continue
			}
			if state != armresources.ProvisioningStateFailed && state != armresources.ProvisioningStateSucceeded {
				id, err := h.notify(context.ctx, &entry.Message)
				if err == nil {
					log.Tracef("notification sent for stage [%s] with id [%s]", entry.StageId, id)
					entry.Sent()
				}
			}
		}
	}

	h.db.Save(&context.Notification)

	if context.Notification.AllSent() {
		context.Notification.Done()
		context.Done()
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
		results[stageId] = armresources.ProvisioningState(*resource.Properties.ProvisioningState)
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
