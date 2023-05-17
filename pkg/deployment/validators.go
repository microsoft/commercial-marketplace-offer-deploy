package deployment

import (
	"context"
	"errors"
	"fmt"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/to"
	log "github.com/sirupsen/logrus"
)

type DryRunValidator interface {
	Validate(input DryRunValidationInput) (*DryRunResponse, error)
}
type DryRunValidationInput struct {
	ctx context.Context
	cred azcore.TokenCredential
	azureDeployment *AzureDeployment
}


type ValidatorFunc func(input DryRunValidationInput) (*DryRunResponse, error)

func (f ValidatorFunc) Validate(input DryRunValidationInput) (*DryRunResponse, error) {
	return f(input)
}

func validateParams(input DryRunValidationInput) (*DryRunResponse, error) {
	response := &DryRunResponse{
		DryRunResult: DryRunResult{
			Status: to.Ptr("ValidateParams - Succeeded"),
		},
	}
	
	template := input.azureDeployment.Template
	requiredParams := getRequiredParams(template["parameters"].(map[string]interface{}))
	params := input.azureDeployment.Params

	var missingRequiredParams []string
	for _, v := range requiredParams {
		if _, ok := params[v]; !ok {
			missingRequiredParams = append(missingRequiredParams, v)
		}
	}

	if len(missingRequiredParams) > 0 {
		var addInfo []*ErrorAdditionalInfo
		for _, v := range missingRequiredParams {
			addInfo = append(addInfo, &ErrorAdditionalInfo{
				Type:    to.Ptr("TemplateViolation"),
				Info: to.Ptr(v + " is missing"),
			})
		}

		response.DryRunResult.Status = to.Ptr("Failed")
		log.Debugf("Missing required parameters: %v", missingRequiredParams)
		errorMesg := fmt.Sprintf("Missing required parameters: %v", missingRequiredParams)
		response.DryRunResult.Error = &DryRunErrorResponse{
			Code:    to.Ptr("Invalid Template"),
			Message: to.Ptr(errorMesg),
			Target:  nil,
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

