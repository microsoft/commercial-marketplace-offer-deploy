package deployment

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore/to"
	"github.com/Azure/azure-sdk-for-go/sdk/azidentity"
	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/resources/armresources"
)

func mapResponse(whatIfResponse *armresources.DeploymentsClientWhatIfResponse) (*DryRunResponse, error) {
	
	dryRunErrorResponse, err := mapError(whatIfResponse.Error)
	if err != nil {
		return nil, err
	}

	log.Printf("Before creation of DryRunResult")
	dryRunResult := DryRunResult{
		Status: whatIfResponse.Status,
		Error: dryRunErrorResponse,
	}

	log.Printf("After creation of DryRunResult")
	dryRunResponse := DryRunResponse{
		DryRunResult: dryRunResult,
		//Code: whatIfResponse.,
	}

	return &dryRunResponse, nil
}

func mapError(armResourceResponse *armresources.ErrorResponse) (*DryRunErrorResponse, error) {
	log.Printf("Inside MapError")
	if armResourceResponse == nil {
		log.Printf("returning nil")
		return nil, nil

	}

	var dryRunErrorDetails []*DryRunErrorResponse
	if armResourceResponse.Details != nil && len(armResourceResponse.Details) > 0 {
		for _, v := range armResourceResponse.Details {
			dryRunError, err := mapError(v)
			if err != nil {
				log.Printf("There was an error mapping an error detail")
				return nil, err
			}
			dryRunErrorDetails = append(dryRunErrorDetails, dryRunError)
		}
	}

	dryRunErrorResponse := DryRunErrorResponse{
		Message: armResourceResponse.Message,
		Code: armResourceResponse.Code,
		Target: armResourceResponse.Target,
		Details: dryRunErrorDetails,
	}

	return &dryRunErrorResponse, nil
}

func DeleteResourceGroup(subscriptionId string, resourceGroupName string) error {
	cred, err := azidentity.NewDefaultAzureCredential(nil)
	if err != nil {
		log.Fatal(err)
	}
	ctx := context.Background()
	resourceGroupClient, err := armresources.NewResourceGroupsClient(subscriptionId, cred, nil)
	if err != nil {
		return err
	}

	pollerResp, err := resourceGroupClient.BeginDelete(ctx, resourceGroupName, nil)
	if err != nil {
		return err
	}

	_, err = pollerResp.PollUntilDone(ctx, nil)
	if err != nil {
		return err
	}
	return nil
}

func createResourceGroup(subscriptionId string, resourceGroupName string, location string) (*armresources.ResourceGroup, error) {
	cred, err := azidentity.NewDefaultAzureCredential(nil)
	if err != nil {
		log.Fatal(err)
	}
	ctx := context.Background()

	resourceGroupClient, err := armresources.NewResourceGroupsClient(subscriptionId, cred, nil)
	if err != nil {
		return nil, err
	}

	resourceGroupResp, err := resourceGroupClient.CreateOrUpdate(
		ctx,
		resourceGroupName,
		armresources.ResourceGroup{
			Location: to.Ptr(location),
		},
		nil)
	if err != nil {
		return nil, err
	}
	return &resourceGroupResp.ResourceGroup, nil
}

func Create(deploymentName, subscriptionId, resourceGroupName string, template, params map[string]interface{}) (*armresources.DeploymentExtended, error) {
	cred, err := azidentity.NewDefaultAzureCredential(nil)
	if err != nil {
		log.Fatal(err)
	}
	ctx := context.Background()
	deploymentsClient, err := armresources.NewDeploymentsClient(subscriptionId, cred, nil)
	if err != nil {
		return nil, err
	}

	log.Printf("About to Create a deployment")

	deploymentPollerResp, err := deploymentsClient.BeginCreateOrUpdate(
		ctx,
		resourceGroupName,
		deploymentName,
		armresources.Deployment{
			Properties: &armresources.DeploymentProperties{
				Template:   template,
				Parameters: params,
				Mode:       to.Ptr(armresources.DeploymentModeIncremental),
			},
		},
		nil)

	if err != nil {
		return nil, fmt.Errorf("cannot create deployment: %v", err)
	}

	resp, err := deploymentPollerResp.PollUntilDone(ctx, nil)
	if err != nil {
		return nil, fmt.Errorf("cannot get the create deployment future respone: %v", err)
	}

	return &resp.DeploymentExtended, nil
}

func readJson(path string) (map[string]interface{}, error) {
	templateFile, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, err
	}

	template := make(map[string]interface{})
	if err := json.Unmarshal(templateFile, &template); err != nil {
		return nil, err
	}

	return template, nil
}
