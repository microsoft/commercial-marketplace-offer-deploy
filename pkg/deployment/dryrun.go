package deployment

import (
	"context"
	"errors"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore/to"
	"github.com/Azure/azure-sdk-for-go/sdk/azidentity"
	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/resources/armresources"
	"github.com/microsoft/commercial-marketplace-offer-deploy/sdk"
	log "github.com/sirupsen/logrus"
)

func whatIfValidator(input DryRunValidationInput) (*sdk.DryRunResponse, error) {
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
		ValidatorFunc(validateParams),
		ValidatorFunc(whatIfValidator),
	}
}

func validate(validators []DryRunValidator, input DryRunValidationInput) (*sdk.DryRunResponse, error) {
	var responses []*sdk.DryRunResponse
	for _, validator := range validators {
		res, err := validator.Validate(input)
		if err != nil {
			return res, err
		}
		if res != nil {
			responses = append(responses, res)
		}
	}

	return aggregateResponses(responses), nil
}

func aggregateResponses(responses []*sdk.DryRunResponse) *sdk.DryRunResponse {
	if responses == nil || len(responses) == 0 {
		return nil
	}
	//todo: aggregate responses
	return responses[0]
}

func DryRun(ctx context.Context, azureDeployment *AzureDeployment) (*sdk.DryRunResponse, error) {
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
