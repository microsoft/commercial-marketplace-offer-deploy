package deployment

import (
	"context"

	"errors"
	"fmt"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/to"
	"github.com/microsoft/commercial-marketplace-offer-deploy/sdk"
	log "github.com/sirupsen/logrus"
)

type DryRunValidator interface {
	Validate(input DryRunValidationInput) (*sdk.DryRunResponse, error)
}
type DryRunValidationInput struct {
	ctx             context.Context
	cred            azcore.TokenCredential
	azureDeployment *AzureDeployment
}

func (i DryRunValidationInput) GetParams() map[string]interface{} {
	return i.azureDeployment.GetParams()
}

func (i DryRunValidationInput) GetTemplateParams() map[string]interface{} {
	return i.azureDeployment.GetTemplateParams()
}

type ValidatorFunc func(input DryRunValidationInput) (*sdk.DryRunResponse, error)

func (f ValidatorFunc) Validate(input DryRunValidationInput) (*sdk.DryRunResponse, error) {
	return f(input)
}

func validateParams(input DryRunValidationInput) (*sdk.DryRunResponse, error) {
	response := &sdk.DryRunResponse{
		DryRunResult: sdk.DryRunResult{
			Status: to.Ptr("ValidateParams - Succeeded"),
		},
	}

	template := input.azureDeployment.Template
	requiredParams := getRequiredParams(template["parameters"].(map[string]interface{}))
	params := input.GetParams()

	var missingRequiredParams []string
	for _, v := range requiredParams {
		if _, ok := params[v]; !ok {
			missingRequiredParams = append(missingRequiredParams, v)
		}
	}

	if len(missingRequiredParams) > 0 {
		var addInfo []*sdk.ErrorAdditionalInfo
		for _, v := range missingRequiredParams {
			addInfo = append(addInfo, &sdk.ErrorAdditionalInfo{
				Type: to.Ptr("TemplateViolation"),
				Info: to.Ptr(v + " is missing"),
			})
		}

		response.DryRunResult.Status = to.Ptr("Failed")
		log.Debugf("Missing required parameters: %v", missingRequiredParams)
		errorMesg := fmt.Sprintf("Missing required parameters: %v", missingRequiredParams)
		response.DryRunResult.Error = &sdk.DryRunErrorResponse{
			Code:           to.Ptr("Invalid Template"),
			Message:        to.Ptr(errorMesg),
			Target:         nil,
			AdditionalInfo: addInfo,
		}
		return response, errors.New(errorMesg)
	}

	return response, nil
}

func getRequiredParams(params map[string]interface{}) []string {
	var requiredParams []string

	for k, v := range params {
		if v != nil {
			paramMap := v.(map[string]interface{})
			if _, ok := paramMap["defaultValue"]; !ok {
				requiredParams = append(requiredParams, k)
			}
		}
	}

	return requiredParams
}
