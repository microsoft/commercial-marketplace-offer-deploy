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
	// Validates dry run input
	//	returns:
	//		- dryRunError: error that will be returned if there were no runtime errors, AND the validator failed. If there were no errors
	//			and the validator passed, then this will be nil
	//		- error: runtime error, so there won't be a dry run error returned by the validator
	//
	Validate(input DryRunValidationInput) (*sdk.DryRunError, error)
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

type ValidatorFunc func(input DryRunValidationInput) (*sdk.DryRunError, error)

func (f ValidatorFunc) Validate(input DryRunValidationInput) (*sdk.DryRunError, error) {
	return f(input)
}

func validateParams(input DryRunValidationInput) (*sdk.DryRunError, error) {
	template := input.azureDeployment.Template

	if _, paramsExist := template["parameters"]; !paramsExist {
		return nil, errors.New("template does not have parameters")
	}

	if _, isCorrectType := template["parameters"].(map[string]interface{}); !isCorrectType {
		return nil, errors.New("template parameters are not of type map[string]interface{}")
	}

	templateParams := template["parameters"].(map[string]interface{})
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

		log.Debugf("Missing required parameters: %v", missingRequiredParams)
		errorMesg := fmt.Sprintf("Missing required parameters: %v", missingRequiredParams)

		return &sdk.DryRunError{
			Code:           to.Ptr("Invalid Template"),
			Message:        to.Ptr(errorMesg),
			Target:         nil,
			AdditionalInfo: addInfo,
		}, nil
	}

	// passed, so return nil for errors
	return nil, nil
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
