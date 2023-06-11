package stage

import (
	"context"
	"errors"
	"time"

	"github.com/Azure/azure-sdk-for-go/sdk/azidentity"
	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/resources/armresources"
	"github.com/microsoft/commercial-marketplace-offer-deploy/internal/model"
	"github.com/microsoft/commercial-marketplace-offer-deploy/internal/model/operation"
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
	ProvisioningState armresources.ProvisioningState `json:"provisioningState"`
}

type DeployStagePollerFactory struct {
}

func NewDeployStagePollerFactory() *DeployStagePollerFactory {
	return &DeployStagePollerFactory{}
}

func (factory *DeployStagePollerFactory) Create(operation *operation.Operation, options *DeployStagePollerOptions) (*DeployStagePoller, error) {
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
		done:              make(chan DeployStagePollerResponse),
		resourceGroupName: deployment.ResourceGroup,
		deploymentName:    factory.getAzureDeploymentName(operation),
	}, nil
}

func (factory *DeployStagePollerFactory) getAzureDeploymentName(operation *operation.Operation) string {
	parameterValue, ok := operation.ParameterValue(model.ParameterKeyAzureDeploymentName)
	if ok {
		if value, ok := parameterValue.(string); ok {
			return value
		}
	}
	return ""
}

func (poller *DeployStagePoller) PollUntilDone(ctx context.Context) (DeployStagePollerResponse, error) {
	for {
		select {
		case <-poller.ticker.C:
			state, err := poller.checkProvisioningState(ctx)
			log.Tracef("provisioning state of stage deployment: %v", state)

			if err != nil {
				log.Errorf("failed to check provisioning state: %v", err)
			}
			if poller.isInCompletedState(state) {
				poller.done <- DeployStagePollerResponse{
					ProvisioningState: state,
				}
			}
		case response := <-poller.done:
			return response, nil
		}
	}
}

func (poller *DeployStagePoller) checkProvisioningState(ctx context.Context) (armresources.ProvisioningState, error) {
	response, err := poller.client.Get(ctx, poller.resourceGroupName, poller.deploymentName, nil)
	if err != nil {
		return "", err
	}
	state := response.DeploymentExtended.Properties.ProvisioningState
	if state == nil {
		return "", errors.New("provisioningState is nil")
	}
	return *state, nil
}

func (poller *DeployStagePoller) isInCompletedState(state armresources.ProvisioningState) bool {
	return state == armresources.ProvisioningStateSucceeded || state == armresources.ProvisioningStateFailed || state == armresources.ProvisioningStateCanceled
}
