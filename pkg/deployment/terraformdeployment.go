package deployment

import (
	"errors"

	"github.com/google/uuid"
	log "github.com/sirupsen/logrus"
)

type TerraformTemplate map[string]interface{}
type TerraformTemplateParameters map[string]interface{}

type TerraformDeployment struct {
	SubscriptionId    string                      `json:"subscriptionId"`
	Location          string                      `json:"location"`
	ResourceGroupName string                      `json:"resourceGroupName"`
	DeploymentName    string                      `json:"deploymentName"`
	DeploymentType    DeploymentType              `json:"deploymentType"`
	Template          TerraformTemplate           `json:"template"`
	Params            TerraformTemplateParameters `json:"templateParams"`
	WorkingDirectory  string                      `json:"workingDirectory"`
	Tags              map[string]*string          `json:"tags"`
	ResumeToken       string                      `json:"resumeToken"`
	OperationId       uuid.UUID                   `json:"operationId"` //the modm operationId that triggered the deployment
}

type BeginTerraformDeploymentResult struct {
	CorrelationID *string     `json:"correlationId" mapstructure:"correlationId"`
	ResumeToken   ResumeToken `json:"resumeToken" mapstructure:"resumeToken"`
}

func (td *TerraformDeployment) logger() *log.Entry {
	return log.WithFields(log.Fields{
		"subscriptionId":    td.SubscriptionId,
		"resourceGroupName": td.ResourceGroupName,
		"deploymentName":    td.DeploymentName,
	})
}

func (ad *TerraformDeployment) GetParameters() map[string]interface{} {
	return getParams(ad.Params)
}

func (ad *TerraformDeployment) GetParametersFromTemplate() map[string]interface{} {
	return getParams(ad.Template)
}

func (ad *TerraformDeployment) GetDeploymentType() DeploymentType {
	return ad.DeploymentType
}

func (td *TerraformDeployment) validate() error {
	if len(td.SubscriptionId) == 0 {
		return errors.New("subscriptionId is not set on azureDeployment input struct")
	}
	if len(td.Location) == 0 {
		return errors.New("location is not set on azureDeployment input struct")
	}
	if len(td.ResourceGroupName) == 0 {
		return errors.New("resourceGroupName is not set on azureDeployment input struct")
	}
	if len(td.ResourceGroupName) == 0 {
		return errors.New("resourceGroupName is not set on azureDeployment input struct")
	}
	if td.Template == nil {
		return errors.New("template is not set on deployment azureDeployment struct")
	}
	// allow params to be empty to support all default params
	return nil
}
