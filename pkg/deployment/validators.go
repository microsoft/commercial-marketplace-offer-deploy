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
	Validate(input DryRunValidationInput) (*sdk.DryRunResult, error)
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

type ValidatorFunc func(input DryRunValidationInput) (*sdk.DryRunResult, error)

func (f ValidatorFunc) Validate(input DryRunValidationInput) (*sdk.DryRunResult, error) {
	return f(input)
}

func validateParams(input DryRunValidationInput) (*sdk.DryRunResult, error) {
	// response := &sdk.DryRunResponse{
	// 	DryRunResult: sdk.DryRunResult{
	// 		Status: to.Ptr(sdk.StatusSuccess.String()),
	// 	},
	// }
	dryRunResult := &sdk.DryRunResult {
		Status: to.Ptr(sdk.StatusSuccess.String()),
	}

	template := input.azureDeployment.Template

	var templateParams map[string]interface{}
	if val, ok := template["parameters"]; ok {
		templateParams, ok = val.(map[string]interface{})
		if !ok {
			// Handle the type assertion failure
			log.Error("error: 'parameters' is not of type map[string]interface{}")
			// This is considered a runtime error, not a validation error
			return nil, errors.New("error: 'parameters' is not of type map[string]interface{}")
		}
	} else {
		// Handle the case where "parameters" key is not present or is nil
		log.Error("Error: 'parameters' key is missing or is nil")
		// Additional error handling if needed
		return nil, errors.New("error: 'parameters' key is missing or is nil")
	}

	requiredParams := getRequiredParams(templateParams)
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

		dryRunResult.Status = to.Ptr("Failed")
		log.Debugf("Missing required parameters: %v", missingRequiredParams)
		errorMesg := fmt.Sprintf("Missing required parameters: %v", missingRequiredParams)
		dryRunResult.Error = &sdk.DryRunErrorResponse{
			Code:           to.Ptr("Invalid Template"),
			Message:        to.Ptr(errorMesg),
			Target:         nil,
			AdditionalInfo: addInfo,
		}
		return dryRunResult, errors.New(errorMesg)
	}

	return dryRunResult, nil
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
