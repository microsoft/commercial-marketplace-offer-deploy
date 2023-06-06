package notification

import (
	"context"
	"fmt"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore/to"
	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/resources/armresources"
	"github.com/microsoft/commercial-marketplace-offer-deploy/internal/hook"
	"github.com/microsoft/commercial-marketplace-offer-deploy/internal/model"
	log "github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

type StageNotificationHandlerFactoryFunc func() (*StageNotificationHandler, error)

type stageNotificationHandlerContext struct {
	ctx          context.Context
	notification *model.StageNotification
	done         chan stageNotificationHandlerResult
}

type stageNotificationHandlerResult struct {
	Id    uint
	Error error
}

type StageNotificationHandler struct {
	db                *gorm.DB
	notify            hook.NotifyFunc
	deploymentsClient *armresources.DeploymentsClient
}

func NewStageNotificationHandler(db *gorm.DB, deploymentsClient *armresources.DeploymentsClient) *StageNotificationHandler {
	return &StageNotificationHandler{
		db:                db,
		deploymentsClient: deploymentsClient,
		notify:            hook.Notify,
	}
}

func (h *StageNotificationHandler) Handle(context *stageNotificationHandlerContext) {
	resources, err := h.getAzureDeploymentResources(context.ctx, context.notification)
	if err != nil {
		context.done <- stageNotificationHandlerResult{
			Id:    context.notification.ID,
			Error: err,
		}
		return
	}

	// TODO: handle the stage notifications until all are sent

	for _, resource := range resources {
		log.Debugf("Handle stage notification for deployment %s", *resource.Name)
	}

	//finally after handling the entire set of notifications to send
	context.done <- stageNotificationHandlerResult{
		Id: context.notification.ID,
	}
}

// get by correlationId
func (h *StageNotificationHandler) getAzureDeploymentResources(ctx context.Context, notification *model.StageNotification) ([]*armresources.DeploymentExtended, error) {
	filter := fmt.Sprintf("correlationId eq '%s'", notification.CorrelationId.String())
	pager := h.deploymentsClient.NewListByResourceGroupPager(notification.ResourceGroupName, &armresources.DeploymentsClientListByResourceGroupOptions{
		Filter: to.Ptr(filter),
	})

	resources := []*armresources.DeploymentExtended{}

	for pager.More() {
		nextResult, err := pager.NextPage(ctx)
		if err != nil {
			return nil, err
		}
		if nextResult.DeploymentListResult.Value != nil {
			resources = append(resources, nextResult.DeploymentListResult.Value...)
		}
	}

	return resources, nil
}
