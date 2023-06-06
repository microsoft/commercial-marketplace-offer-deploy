package template

import (
	"strconv"

	"github.com/google/uuid"
	"github.com/microsoft/commercial-marketplace-offer-deploy/pkg/deployment"
)

const deploymentResourceTypeName = "Microsoft.Resources/deployments"
const DefaultNoRetries = 0

type ArmTemplate struct {
	Resources []ArmTemplateResource `mapstructure:"resources"`
}

type ArmTemplateResource struct {
	Name string            `mapstructure:"name"`
	Type string            `mapstructure:"type"`
	Tags map[string]string `mapstructure:"tags"`
}

func (resource *ArmTemplateResource) IsDeploymentResourceType() bool {
	return resource.Type == deploymentResourceTypeName
}

func (resource *ArmTemplateResource) GetId() uuid.UUID {
	defaultValue := uuid.NewString()
	value := resource.GetTagValue(deployment.LookupTagKeyId, defaultValue)
	return uuid.MustParse(value)
}

func (resource *ArmTemplateResource) GetName() string {
	defaultValue := resource.Name
	return resource.GetTagValue(deployment.LookupTagKeyName, defaultValue)
}

// gets a stage's default retry value via modm tag. If the value is not an integer, the default is used.
//
//	default: DefaultNoRetries
func (resource *ArmTemplateResource) GetRetries() uint {
	defaultValue := DefaultNoRetries
	value := resource.GetTagValue(deployment.LookupTagKeyRetry, strconv.Itoa(defaultValue))

	if intValue, err := strconv.ParseInt(value, 10, 32); err == nil {
		return uint(intValue)
	}
	return uint(defaultValue)
}

func (resource *ArmTemplateResource) GetTagValue(key deployment.LookupTagKey, defaultValue string) string {
	if resource.Tags != nil {
		if value, exists := resource.Tags[string(key)]; exists {
			return value
		}
	}
	return defaultValue
}
