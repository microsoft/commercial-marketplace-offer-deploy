package deployment

import (
	"context"
	"log"
<<<<<<< HEAD
=======

>>>>>>> main
	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/to"
	"github.com/Azure/azure-sdk-for-go/sdk/azidentity"
	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/resources/armresources"
)

type DryRunResponse struct {
	DryRunResult
}

type DryRunResult struct {
	Status *string `json:"code,omitempty" azure:"ro"`
	//Message *string `json:"message,omitempty" azure:"ro"`
	//Target *string `json:"target,omitempty" azure:"ro"`
	Error *DryRunErrorResponse
}

type DryRunErrorResponse struct {
	// READ-ONLY; The error additional info.
	AdditionalInfo []*ErrorAdditionalInfo `json:"additionalInfo,omitempty" azure:"ro"`

	// READ-ONLY; The error code.
	Code *string `json:"code,omitempty" azure:"ro"`

	// READ-ONLY; The error message.
	Message *string `json:"message,omitempty" azure:"ro"`

	// READ-ONLY; The error target.
	Target *string `json:"target,omitempty" azure:"ro"`

	// READ-ONLY; The error details.
	Details []*DryRunErrorResponse `json:"details,omitempty" azure:"ro"`
}

<<<<<<< HEAD
func whatIfValidator(azureDeployment *AzureDeployment) *DryRunResponse {
=======
func DryRun(azureDeployment *AzureDeployment) *DryRunResponse {
>>>>>>> main
	err := azureDeployment.validate()
	if err != nil {
		log.Fatal(err)
	}
<<<<<<< HEAD

	cred, err := azidentity.NewDefaultAzureCredential(nil)
	
=======
	cred, err := azidentity.NewDefaultAzureCredential(nil)
>>>>>>> main
	if err != nil {
		log.Fatal(err)
	}
	ctx := context.Background()

	whatIfResult, err := whatIfDeployment(ctx, cred, azureDeployment)
	if err != nil {
		log.Fatal(err)
	}

	dryResponse, err := mapResponse(whatIfResult)
	if err != nil {
		log.Fatal(err)
	}

	return dryResponse
}

func loadValidators() []DryRunValidator {
	//todo: load validators from config
	return []DryRunValidator{
		WhatIfValidatorFunc(whatIfValidator),
	}
}


func validate(validators []DryRunValidator, azureDeployment *AzureDeployment) *DryRunResponse {
	var responses []*DryRunResponse
	for _, validator := range validators {
		res := validator.Validate(azureDeployment)
		if res != nil {
			responses = append(responses, res)
		}
	}
	
	return aggregateResponses(responses)
}

func aggregateResponses(responses []*DryRunResponse) *DryRunResponse {
	if responses == nil || len(responses) == 0 {
		return nil
	}
	//todo: aggregate responses
	return responses[0]
}

func DryRun(azureDeployment *AzureDeployment) *DryRunResponse {
	validators := loadValidators()
	return validate(validators, azureDeployment)
}

func whatIfDeployment(ctx context.Context, cred azcore.TokenCredential, azureDeployment *AzureDeployment) (*armresources.DeploymentsClientWhatIfResponse, error) {
	deploymentsClient, err := armresources.NewDeploymentsClient(azureDeployment.subscriptionId, cred, nil)
	if err != nil {
		return nil, err
	}

	pollerResp, err := deploymentsClient.BeginWhatIf(
		ctx,
		azureDeployment.resourceGroupName,
		azureDeployment.deploymentName,
		armresources.DeploymentWhatIf{
			Properties: &armresources.DeploymentWhatIfProperties{
				Template:   azureDeployment.template,
				Parameters: azureDeployment.params,
				Mode:       to.Ptr(armresources.DeploymentModeIncremental),
			},
		},
		nil)

	if err != nil {
		return nil, err
	}

	resp, err := pollerResp.PollUntilDone(ctx, nil)
	if err != nil {
		return nil, err
	}

	return &resp, nil
}
