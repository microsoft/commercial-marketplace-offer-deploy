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

type StageNotificationService struct {
	ctx               context.Context
	deploymentsClient *armresources.DeploymentsClient
	notify            hook.NotifyFunc
	pump              *StageNotificationPump
	notifications     map[uint]*model.StageNotification
	db                *gorm.DB
	log               *log.Entry
}

func NewStageNotificationService() *StageNotificationService {
	return &StageNotificationService{
		log: log.WithFields(log.Fields{}),
	}
}

// stub out hosting.Service interface on StageNotificationService
func (s *StageNotificationService) Start() {
	s.pump.Start()
}

func (s *StageNotificationService) Stop() {
	s.pump.Stop()
}

func (s *StageNotificationService) GetName() string {
	return ""
}

// get by correlationId
func (s *StageNotificationService) getAzureDeploymentResources(notification *model.StageNotification) []*armresources.DeploymentExtended {
	filter := fmt.Sprintf("correlationId eq '%s'", notification.CorrelationId.String())
	response := s.deploymentsClient.NewListByResourceGroupPager(notification.ResourceGroupName, &armresources.DeploymentsClientListByResourceGroupOptions{
		Filter: to.Ptr(filter),
	})

	//TODO: add pager loop and collect all deployments
	response.NextPage(s.ctx)

	return nil
}
