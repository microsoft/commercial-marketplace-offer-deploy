package operation

import (
	"context"
	"errors"
	"time"

	"github.com/Azure/azure-sdk-for-go/sdk/azidentity"
	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/resources/armresources"
	"github.com/google/uuid"
	"github.com/microsoft/commercial-marketplace-offer-deploy/pkg/deployment"
	log "github.com/sirupsen/logrus"
)

type FinderResponse struct {
	Name string
}
type AzureDeploymentNameFinder struct {
	client            *armresources.DeploymentsClient
	operationId       uuid.UUID
	ticker			  *time.Ticker
	done 			  chan FinderResponse
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
		ticker:            time.NewTicker(10 * time.Second),
		done:              make(chan FinderResponse, 0),
		resourceGroupName: deployment.ResourceGroup,
		operationId:       operation.ID,
	}, nil
}

func (finder *AzureDeploymentNameFinder) FindUntilDone(ctx context.Context) (string, error) {
	for {
		select {
		case <-finder.ticker.C:
			log.Tracef("Finding deployment name for operationId: %s", finder.operationId)
			name, err := finder.getName(ctx)
			if err != nil {
				log.Errorf("Failed to find deployment name for operationId: %s", finder.operationId)
			}
			if len(name) > 0 {
				finder.done <- FinderResponse{
					Name: name,
				}
				return name, nil
			}
		case response := <-finder.done:
			log.Tracef("Found deployment name for operationId: %s", finder.operationId)
			return response.Name, nil
		case <-ctx.Done():
			return "", ctx.Err()
		}
	}
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

				if value, ok := (item.Tags)[string(deployment.LookupTagKeyOperationId)]; ok {
					id, err := uuid.Parse(*value)
					if err != nil {
						log.Warnf("Failed to parse operationId from deployment tags: %s", err.Error())
						continue
					}

					if id != finder.operationId {
						continue
					}

					name = *item.Name
					log.WithFields(log.Fields{
						"operationId":         finder.operationId,
						"azureDeploymentName": name,
					}).Trace("Found deployment by operationId")
					break
				}
			}
		}
	}

	return name, nil
}
