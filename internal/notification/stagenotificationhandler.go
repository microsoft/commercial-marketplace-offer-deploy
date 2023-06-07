package notification

import (
	"context"

	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/resources/armresources"
	"github.com/google/uuid"
	"github.com/microsoft/commercial-marketplace-offer-deploy/internal/hook"
	"github.com/microsoft/commercial-marketplace-offer-deploy/internal/model"
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

func NewStageNotificationHandler(db *gorm.DB, deploymentsClient *armresources.DeploymentsClient) NotificationHandler[model.StageNotification] {
	return &stageNotificationHandler{
		db:                db,
		deploymentsClient: deploymentsClient,
		notify:            hook.Notify,
	}
}

func (h *stageNotificationHandler) Handle(context *NotificationHandlerContext[model.StageNotification]) {
	resources, err := h.getAzureDeploymentResources(context.ctx, context.Notification)
	if err != nil {
		context.Done(NotificationHandlerResult[model.StageNotification]{
			Notification: context.Notification,
			Error:        err,
		})
		return
	}

	// TODO: handle the stage notifications until all are sent

	for _, resource := range resources {
		log.Debugf("Handle stage notification for deployment %s", *resource.Name)
	}

	context.Done(NotificationHandlerResult[model.StageNotification]{
		Notification: context.Notification,
	})
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
