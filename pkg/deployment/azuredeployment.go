package deployment

import (
	"errors"
	"time"

	"github.com/google/uuid"
	log "github.com/sirupsen/logrus"
)

const deploymentResourceTypeName = "Microsoft.Resources/deployments"

type AzureTemplate map[string]interface{}
type AzureTemplateParameters map[string]interface{}

type AzureDeployment struct {
	SubscriptionId    string                  `json:"subscriptionId"`
	Location          string                  `json:"location"`
	ResourceGroupName string                  `json:"resourceGroupName"`
	DeploymentName    string                  `json:"deploymentName"`
	DeploymentType    DeploymentType          `json:"deploymentType"`
	Template          AzureTemplate           `json:"template"`
	Params            AzureTemplateParameters `json:"templateParams"`
	Tags              map[string]*string      `json:"tags"`
	ResumeToken       string                  `json:"resumeToken"`
	OperationId       uuid.UUID               `json:"operationId"` //the modm operationId that triggered the deployment
}

func (ad *AzureDeployment) logger() *log.Entry {
	return log.WithFields(log.Fields{
		"subscriptionId":    ad.SubscriptionId,
		"resourceGroupName": ad.ResourceGroupName,
		"deploymentName":    ad.DeploymentName,
	})
}

// gets the template
func (ad *AzureDeployment) GetTemplate() map[string]interface{} {
	if ad.Template == nil {
		return nil
	}
	operationIdKey := string(LookupTagKeyOperationId)
	template := ad.Template

	// set the operationId on every nested template at level 1 of the parent template
	if resourcesEntry, ok := template["resources"]; ok {
		if resources, ok := resourcesEntry.([]any); ok {
			if len(resources) > 0 {
				for _, resourceEntry := range resources {
					if resourceMap, ok := resourceEntry.(map[string]any); ok {
						if isDeploymentResourceType(resourceMap) {
							if tagsEntry, ok := resourceMap["tags"]; ok {
								if tagsMap, ok := tagsEntry.(map[string]any); ok {
									tagsMap[operationIdKey] = ad.OperationId.String()
								}
							}
						}
					}
				}
			}
		}
	}
	return template //return the modified copy
}

func (ad *AzureDeployment) GetParameters() map[string]interface{} {
	return getParams(ad.Params)
}

func (ad *AzureDeployment) GetParametersFromTemplate() map[string]interface{} {
	return getParams(ad.Template)
}

func (ad *AzureDeployment) GetDeploymentType() DeploymentType {
	return ad.DeploymentType
}

func (azureDeployment *AzureDeployment) validate() error {
	if len(azureDeployment.SubscriptionId) == 0 {
		return errors.New("subscriptionId is not set on azureDeployment input struct")
	}
	if len(azureDeployment.Location) == 0 {
		return errors.New("location is not set on azureDeployment input struct")
	}
	if len(azureDeployment.ResourceGroupName) == 0 {
		return errors.New("resourceGroupName is not set on azureDeployment input struct")
	}
	if len(azureDeployment.ResourceGroupName) == 0 {
		return errors.New("resourceGroupName is not set on azureDeployment input struct")
	}
	if azureDeployment.Template == nil {
		return errors.New("template is not set on deployment azureDeployment struct")
	}
	// allow params to be empty to support all default params
	return nil
}

type AzureRedeployment struct {
	SubscriptionId    string             `json:"subscriptionId"`
	Location          string             `json:"location"`
	ResourceGroupName string             `json:"resourceGroupName"`
	DeploymentName    string             `json:"deploymentName"`
	Tags              map[string]*string `json:"tags"`
}

type AzureDeploymentResult struct {
	ID                string                 `json:"id"`
	CorrelationID     string                 `json:"correlationId"`
	Duration          string                 `json:"duration"`
	Timestamp         time.Time              `json:"timestamp"`
	ProvisioningState string                 `json:"provisioningState"`
	Outputs           map[string]interface{} `json:"outputs"`
	Status            ExecutionStatus
}

type AzureCancelDeployment struct {
	SubscriptionId    string `json:"subscriptionId"`
	Location          string `json:"location"`
	ResourceGroupName string `json:"resourceGroupName"`
	DeploymentName    string `json:"deploymentName"`
}

type AzureCancelDeploymentResult struct {
	CancelSubmitted bool `json:"cancelSubmitted"`
}

type BeginAzureDeploymentResult struct {
	CorrelationID *string     `json:"correlationId" mapstructure:"correlationId"`
	ResumeToken   ResumeToken `json:"resumeToken" mapstructure:"resumeToken"`
}

// the resume token to resume waiting for an azure operation
type ResumeToken struct {
	SubscriptionId string `json:"subscriptionId" mapstructure:"subscriptionId"`
	Value          string `json:"value" mapstructure:"value"`
}

func isDeploymentResourceType(resourceMap map[string]any) bool {
	if resourceMap == nil {
		return false
	}
	if typeEntry, ok := resourceMap["type"]; ok {
		if typeValue, ok := typeEntry.(string); ok {
			return typeValue == deploymentResourceTypeName
		}
	}
	return false
}
