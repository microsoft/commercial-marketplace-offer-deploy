package deployment

import (
	"context"
	"errors"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore/to"
	"github.com/Azure/azure-sdk-for-go/sdk/azidentity"
	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/resources/armresources"
	log "github.com/sirupsen/logrus"
)

type DryRunResponse struct {
	DryRunResult
}

type DryRunResult struct {
	Status *string `json:"status,omitempty" azure:"ro"`
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

func whatIfValidator(input DryRunValidationInput) (*DryRunResponse, error) {
	if input.azureDeployment == nil {
		log.Error(errors.New("azureDeployment is not set on input struct"))
	}
	err := input.azureDeployment.validate()
	if err != nil {
		return nil, err
	}

	whatIfResult, err := whatIfDeployment(input)
	if err != nil {
		return nil, err
	}

	dryResponse, err := mapResponse(whatIfResult)
	if err != nil {
		return nil, err
	}

	return dryResponse, nil
}

func loadValidators() []DryRunValidator {
	//todo: load validators from config
	return []DryRunValidator{
		WhatIfValidatorFunc(whatIfValidator),
	}
}

func validate(validators []DryRunValidator, input DryRunValidationInput) (*DryRunResponse, error) {
	var responses []*DryRunResponse
	for _, validator := range validators {
		res, err := validator.Validate(input)
		if err != nil {
			return nil, err
		}
		if res != nil {
			responses = append(responses, res)
		}
	}

	return aggregateResponses(responses), nil
}

func aggregateResponses(responses []*DryRunResponse) *DryRunResponse {
	if responses == nil || len(responses) == 0 {
		return nil
	}
	//todo: aggregate responses
	return responses[0]
}

func DryRun(ctx context.Context, azureDeployment *AzureDeployment) (*DryRunResponse, error) {
	log.Debug("Inside DryRun in pkg/deployment/dryrun.go")
	cred, err := azidentity.NewDefaultAzureCredential(nil)
	if err != nil {
		return nil, err
	}
	validators := loadValidators()
	input := DryRunValidationInput{
		ctx:             ctx,
		cred:            cred,
		azureDeployment: azureDeployment,
	}
	log.Debug("About to call DryRun validator")
	return validate(validators, input)
}

func whatIfDeployment(input DryRunValidationInput) (*armresources.DeploymentsClientWhatIfResponse, error) {
	log.Debug("Inside whatIfDeployment in pkg/deployment/dryrun.go")
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

	deploymentsClient, err := armresources.NewDeploymentsClient(azureDeployment.SubscriptionId, cred, nil)
	if err != nil {
		return nil, err
	}

	var templateParams map[string]interface{}
	if azureDeployment.Params != nil {
		if p, ok := azureDeployment.Params["parameters"]; ok {
			templateParams = p.(map[string]interface{})
		} else {
			templateParams = azureDeployment.Params
		}
	}

	log.Debugf("About to call whatIf with templateParams of - %v", templateParams)

	pollerResp, err := deploymentsClient.BeginWhatIf(
		ctx,
		azureDeployment.ResourceGroupName,
		azureDeployment.DeploymentName,
		armresources.DeploymentWhatIf{
			Properties: &armresources.DeploymentWhatIfProperties{
				Template:   azureDeployment.Template,
				Parameters: templateParams,
				Mode:       to.Ptr(armresources.DeploymentModeIncremental),
			},
		},
		nil)

	if err != nil {
		return nil, err
	}

	log.Debug("Got the whatIf response")

	resp, err := pollerResp.PollUntilDone(ctx, nil)

	if err != nil {
		return nil, err
	}
	log.Debugf("whatIf response - %v", resp)

	return &resp, nil
}
