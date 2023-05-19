package deployment

import (
	"context"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore/to"
	"github.com/Azure/azure-sdk-for-go/sdk/azidentity"
	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/resources/armresources"
	"github.com/microsoft/commercial-marketplace-offer-deploy/internal/utils"
	"github.com/microsoft/commercial-marketplace-offer-deploy/sdk"
	log "github.com/sirupsen/logrus"
)

type DeploymentType int64
type Template map[string]interface{}
type TemplateParams map[string]interface{}

// create deployment function
type CreateDeployment func(ctx context.Context, dep AzureDeployment) (*AzureDeploymentResult, error)

const (
	AzureResourceManager DeploymentType = iota
	Terraform
)

func mapResult(whatIfResponse *armresources.DeploymentsClientWhatIfResponse) (*sdk.DryRunError, error) {
	dryRunError, err := mapError(whatIfResponse.Error)
	if err != nil {
		return nil, err
	}
	return dryRunError, nil
}

func mapError(armResourceResponse *armresources.ErrorResponse) (*sdk.DryRunError, error) {
	if armResourceResponse == nil {
		log.Debug("returning nil")
		return nil, nil
	}

	var dryRunErrorDetails []*sdk.DryRunError
	if armResourceResponse.Details != nil && len(armResourceResponse.Details) > 0 {
		for _, v := range armResourceResponse.Details {
			dryRunError, err := mapError(v)
			if err != nil {
				log.Error("There was an error mapping an error detail")
				return nil, err
			}
			dryRunErrorDetails = append(dryRunErrorDetails, dryRunError)
		}
	}

	var errorAdditionalInfo []*sdk.ErrorAdditionalInfo
	if armResourceResponse.AdditionalInfo != nil && len(armResourceResponse.AdditionalInfo) > 0 {
		for _, v := range armResourceResponse.AdditionalInfo {
			errAddInfo := &sdk.ErrorAdditionalInfo{Info: v.Info, Type: v.Type}
			errorAdditionalInfo = append(errorAdditionalInfo, errAddInfo)
		}
	}

	dryRunErrorResponse := sdk.DryRunError{
		Message:        armResourceResponse.Message,
		Code:           armResourceResponse.Code,
		Target:         armResourceResponse.Target,
		Details:        dryRunErrorDetails,
		AdditionalInfo: errorAdditionalInfo,
	}

	return &dryRunErrorResponse, nil
}

func DeleteResourceGroup(ctx context.Context, subscriptionId string, resourceGroupName string) error {
	cred, err := azidentity.NewDefaultAzureCredential(nil)
	if err != nil {
		log.Error(err)
	}

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
		log.Error(err)
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

func Create(ctx context.Context, dep AzureDeployment) (*AzureDeploymentResult, error) {
	log.Debug("Inside Create")
	deployer := CreateNewDeployer(dep.DeploymentType)
	return deployer.Deploy(ctx, &dep)
}

func Redeploy(ctx context.Context, dep AzureRedeployment) (*AzureDeploymentResult, error) {
	log.Debug("Inside Redeploy")
	deploymentType := 0
	deployer := CreateNewDeployer(DeploymentType(deploymentType))
	return deployer.Redeploy(ctx, &dep)
}

func readJson(path string) (map[string]interface{}, error) {
	return utils.ReadJson(path)
}
