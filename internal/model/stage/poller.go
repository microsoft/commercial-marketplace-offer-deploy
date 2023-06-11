package stage

import (
	"context"
	"errors"
	"time"

	"github.com/Azure/azure-sdk-for-go/sdk/azidentity"
	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/resources/armresources"
	"github.com/microsoft/commercial-marketplace-offer-deploy/internal/model/operation"
	"github.com/microsoft/commercial-marketplace-offer-deploy/sdk"
	log "github.com/sirupsen/logrus"
)

const DefaultFrequency = 10 * time.Second

type DeployStagePollerOptions struct {
	Frequency time.Duration
}

type DeployStagePoller struct {
	client            *armresources.DeploymentsClient
	ticker            *time.Ticker
	done              chan DeployStagePollerResponse
	resourceGroupName string
	deploymentName    string
}

type DeployStagePollerResponse struct {
	Status     sdk.Status                       `json:"status"`
	Deployment *armresources.DeploymentExtended `json:"deployment"`
}

type DeployStagePollerFactory struct {
}

func NewDeployStagePollerFactory() *DeployStagePollerFactory {
	return &DeployStagePollerFactory{}
}

func (factory *DeployStagePollerFactory) Create(operation *operation.Operation, azureDeploymentName string, options *DeployStagePollerOptions) (*DeployStagePoller, error) {
	credential, err := azidentity.NewDefaultAzureCredential(nil)
	if err != nil {
		return nil, err
	}

	deployment := operation.Deployment()
	if deployment == nil {
		return nil, errors.New("deployment is nil. failed to create DeployStagePoller")
	}

	client, err := armresources.NewDeploymentsClient(deployment.SubscriptionId, credential, nil)
	if err != nil {
		return nil, err
	}

	duration := DefaultFrequency
	if options != nil {
		duration = options.Frequency
	}

	return &DeployStagePoller{
		client:            client,
		ticker:            time.NewTicker(duration),
		done:              make(chan DeployStagePollerResponse, 1),
		resourceGroupName: deployment.ResourceGroup,
		deploymentName:    azureDeploymentName,
	}, nil
}

func (poller *DeployStagePoller) PollUntilDone(ctx context.Context) (DeployStagePollerResponse, error) {
	for {
		select {
		case <-poller.ticker.C:
			log.Tracef("checking provisioning state of stage deployment [%s]", poller.deploymentName)

			deployment, err := poller.checkProvisioningState(ctx)
			log.Tracef("provisioning state of stage deployment: %s", *deployment.Properties.ProvisioningState)

			if err != nil {
				log.Errorf("failed to check provisioning state: %v", err)
			}
			if poller.isInCompletedState(*deployment.Properties.ProvisioningState) {
				poller.ticker.Stop()
				poller.done <- DeployStagePollerResponse{
					Status:     poller.mapProvisioningStateToStatus(*deployment.Properties.ProvisioningState),
					Deployment: deployment,
				}
			}
		case response := <-poller.done:
			log.Tracef("poller is done [%s]", response.Status)
			return response, nil
		}
	}
}

func (poller *DeployStagePoller) checkProvisioningState(ctx context.Context) (*armresources.DeploymentExtended, error) {
	response, err := poller.client.Get(ctx, poller.resourceGroupName, poller.deploymentName, nil)
	if err != nil {
		return nil, err
	}
	state := response.DeploymentExtended.Properties.ProvisioningState
	if state == nil {
		return nil, errors.New("provisioningState is nil")
	}
	return &response.DeploymentExtended, nil
}

func (poller *DeployStagePoller) isInCompletedState(state armresources.ProvisioningState) bool {
	return state == armresources.ProvisioningStateSucceeded || state == armresources.ProvisioningStateFailed || state == armresources.ProvisioningStateCanceled
}

func (poller *DeployStagePoller) mapProvisioningStateToStatus(state armresources.ProvisioningState) sdk.Status {
	switch state {
	case armresources.ProvisioningStateSucceeded:
		return sdk.StatusSuccess
	case armresources.ProvisioningStateFailed:
		return sdk.StatusFailed
	case armresources.ProvisioningStateCanceled:
		return sdk.StatusCanceled
	default:
		return sdk.StatusUnknown
	}
}
