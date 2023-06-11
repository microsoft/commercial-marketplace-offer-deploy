package operation

import (
	"context"
	"errors"

	"github.com/Azure/azure-sdk-for-go/sdk/azidentity"
	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/resources/armresources"
	"github.com/google/uuid"
	"github.com/microsoft/commercial-marketplace-offer-deploy/pkg/deployment"
)

type AzureDeploymentNameFinder struct {
	client            *armresources.DeploymentsClient
	operationId       uuid.UUID
	resourceGroupName string
}

func NewAzureDeploymentNameFinder(operation *Operation) (*AzureDeploymentNameFinder, error) {
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
	return &AzureDeploymentNameFinder{
		client:            client,
		resourceGroupName: deployment.ResourceGroup,
		operationId:       operation.ID,
	}, nil
}

func (finder *AzureDeploymentNameFinder) Find(ctx context.Context) (string, error) {
	return finder.getName(ctx)
}

// get by correlationId
func (finder *AzureDeploymentNameFinder) getName(ctx context.Context) (string, error) {
	pager := finder.client.NewListByResourceGroupPager(finder.resourceGroupName, nil)

	name := ""

	for pager.More() {
		nextResult, err := pager.NextPage(ctx)
		if err != nil {
			return name, err
		}
		if nextResult.DeploymentListResult.Value != nil {
			for _, item := range nextResult.DeploymentListResult.Value {
				if item.Tags == nil {
					continue
				}

				if _, ok := (item.Tags)[string(deployment.LookupTagKeyOperationId)]; ok {
					name = *item.Name
					break
				}
			}
		}
	}
	return name, nil
}
