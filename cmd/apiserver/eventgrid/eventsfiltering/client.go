package eventsfiltering

import (
	"context"
	"fmt"
	"strings"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/arm"
	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/resources/armresources"
	log "github.com/sirupsen/logrus"
)

type AzureResourceClient interface {
	Get(ctx context.Context, resourceId *arm.ResourceID) (*armresources.GenericResource, error)
	GetDeployment(ctx context.Context, resourceId *arm.ResourceID) (*armresources.DeploymentExtended, error)
}

type azureResourceClient struct {
	resourcesClient   *armresources.Client
	deploymentsClient *armresources.DeploymentsClient
	providersClient   *armresources.ProvidersClient
}

func NewAzureResourceClient(subscriptionId string, credential azcore.TokenCredential) (AzureResourceClient, error) {
	resourcesClient, err := armresources.NewClient(subscriptionId, credential, nil)
	if err != nil {
		return nil, err
	}

	deploymentsClient, err := armresources.NewDeploymentsClient(subscriptionId, credential, nil)
	if err != nil {
		return nil, err
	}

	providerClient, err := armresources.NewProvidersClient(subscriptionId, credential, nil)
	if err != nil {
		return nil, err
	}

	return &azureResourceClient{
		resourcesClient:   resourcesClient,
		deploymentsClient: deploymentsClient,
		providersClient:   providerClient,
	}, nil
}

func (c *azureResourceClient) Get(ctx context.Context, resourceId *arm.ResourceID) (*armresources.GenericResource, error) {
	apiVersion, err := c.resolveApiVersion(ctx, resourceId)
	if err != nil {
		log.Errorf("err: %v", err)
		return nil, err
	}

	response, err := c.resourcesClient.GetByID(ctx, resourceId.String(), apiVersion, nil)
	if err != nil {
		log.Errorf("failed to get associated resource: %s, err: %v", resourceId.String(), err)
		return nil, err
	}

	return &response.GenericResource, nil
}

func (c *azureResourceClient) GetDeployment(ctx context.Context, resourceId *arm.ResourceID) (*armresources.DeploymentExtended, error) {
	response, err := c.deploymentsClient.Get(ctx, resourceId.Parent.ResourceGroupName, resourceId.Name, nil)
	if err != nil {
		log.Errorf("failed to get deployment resource: %s, err: %v", resourceId.String(), err)
		return nil, err
	}
	return &response.DeploymentExtended, nil
}

// port from Python code from Azure CLI
// reference: https://github.com/Azure/azure-cli/blob/dev/src/azure-cli/azure/cli/command_modules/resource/custom.py

func (c *azureResourceClient) resolveApiVersion(ctx context.Context, resourceId *arm.ResourceID) (string, error) {
	defaultApiVersion := "2021-04-01"

	response, err := c.providersClient.Get(ctx, resourceId.ResourceType.Namespace, nil)
	if err != nil {
		return defaultApiVersion, err
	}
	for _, resourceType := range response.ResourceTypes {
		isResourceTypeMatch := strings.EqualFold(*resourceType.ResourceType, resourceId.ResourceType.Type)
		if isResourceTypeMatch {
			if len(resourceType.APIVersions) > 0 {
				apiVersion := *resourceType.APIVersions[0]
				log.Tracef("resolved api version: %s for resource: %s", apiVersion, resourceId.String())
				return apiVersion, nil
			}
		}
	}
	return defaultApiVersion, fmt.Errorf("failed to resolve api version for resource: %s", resourceId.String())
}
