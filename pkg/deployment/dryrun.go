package deployment

import (
	"context"
	"errors"
	"log"
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

func whatIfValidator(input DryRunValidationInput) *DryRunResponse {
	if input.azureDeployment == nil {
		log.Fatal(errors.New("azureDeployment is not set on input struct"))
	}
	err := input.azureDeployment.validate()
	if err != nil {
		log.Fatal(err)
	}

	whatIfResult, err := whatIfDeployment(input)
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

func validate(validators []DryRunValidator, input DryRunValidationInput) *DryRunResponse {
	var responses []*DryRunResponse
	for _, validator := range validators {
		res := validator.Validate(input)
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
	cred, err := azidentity.NewDefaultAzureCredential(nil)
	if err != nil {
		return nil
	}
	validators := loadValidators()
	input := DryRunValidationInput{
		ctx : context.Background(),
		cred : cred,
		azureDeployment: azureDeployment,
	}
	return validate(validators, input)
}

func whatIfDeployment(input DryRunValidationInput) (*armresources.DeploymentsClientWhatIfResponse, error) {
	if input.azureDeployment == nil {
		return nil, errors.New("azureDeployment is nil")
	}
	azureDeployment := input.azureDeployment
	
	if input.cred == nil {
		return nil, errors.New("credential is nil")
	}
	cred := input.cred

	if input.ctx == nil {
		return nil, errors.New("context is nil")
	}
	ctx := input.ctx
	
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
