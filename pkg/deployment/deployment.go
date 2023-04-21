package deployment

import (
	"context"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore/to"
	"github.com/Azure/azure-sdk-for-go/sdk/azidentity"
	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/resources/armresources"
	"github.com/microsoft/commercial-marketplace-offer-deploy/internal/utils"
	log "github.com/sirupsen/logrus"
)

type DeploymentType int64
type Template map[string]interface{}
type TemplateParams map[string]interface{}

const (
	AzureResourceManager DeploymentType = iota
	Terraform
)

func mapResponse(whatIfResponse *armresources.DeploymentsClientWhatIfResponse) (*DryRunResponse, error) {
	dryRunErrorResponse, err := mapError(whatIfResponse.Error)
	if err != nil {
		return nil, err
	}

	log.Printf("Before creation of DryRunResult")
	dryRunResult := DryRunResult{
		Status: whatIfResponse.Status,
		Error:  dryRunErrorResponse,
	}

	log.Printf("After creation of DryRunResult")
	dryRunResponse := DryRunResponse{
		DryRunResult: dryRunResult,
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

	var errorAdditionalInfo []*ErrorAdditionalInfo
	if armResourceResponse.AdditionalInfo != nil && len(armResourceResponse.AdditionalInfo) > 0 {
		for _, v := range armResourceResponse.AdditionalInfo {
			errAddInfo := &ErrorAdditionalInfo{Info: v.Info, Type: v.Type}
			errorAdditionalInfo = append(errorAdditionalInfo, errAddInfo)
		}
	}

	dryRunErrorResponse := DryRunErrorResponse{
		Message:        armResourceResponse.Message,
		Code:           armResourceResponse.Code,
		Target:         armResourceResponse.Target,
		Details:        dryRunErrorDetails,
		AdditionalInfo: errorAdditionalInfo,
	}

	return &dryRunErrorResponse, nil
}

func DeleteResourceGroup(subscriptionId string, resourceGroupName string) error {
	cred, err := azidentity.NewDefaultAzureCredential(nil)
	if err != nil {
		log.Print(err)
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
		log.Print(err)
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

func Create(dep AzureDeployment) (*AzureDeploymentResult, error) {
	log.Println("Inside Create")
	deployer := CreateNewDeployer(dep)
	return deployer.Deploy(&dep)
}

func readJson(path string) (map[string]interface{}, error) {
	return utils.ReadJson(path)
}
