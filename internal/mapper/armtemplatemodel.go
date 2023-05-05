package mapper

import (
	"strconv"

	"github.com/google/uuid"
	"github.com/microsoft/commercial-marketplace-offer-deploy/pkg/deployment"
)

const deploymentResourceTypeName = "Microsoft.Resources/deployments"
const DefaultNoRetries = 0

type armTemplateResource struct {
	Name string            `mapstructure:"name"`
	Type string            `mapstructure:"type"`
	Tags map[string]string `mapstructure:"tags"`
}

type armTemplate struct {
	Resources []armTemplateResource `mapstructure:"resources"`
}

func (resource *armTemplateResource) isDeploymentResourceType() bool {
	return resource.Type == deploymentResourceTypeName
}

func (resource *armTemplateResource) getId() uuid.UUID {
	defaultValue := uuid.NewString()
	value := resource.getTagValue(deployment.LookupTagKeyId, defaultValue)
	return uuid.MustParse(value)
}

func (resource *armTemplateResource) getName() string {
	defaultValue := resource.Name
	return resource.getTagValue(deployment.LookupTagKeyName, defaultValue)
}

// gets a stage's default retry value via modm tag. If the value is not an integer, the default is used.
//
//	default: DefaultNoRetries
func (resource *armTemplateResource) getRetries() uint {
	defaultValue := DefaultNoRetries
	value := resource.getTagValue(deployment.LookupTagKeyRetry, strconv.Itoa(defaultValue))

	if intValue, err := strconv.ParseInt(value, 10, 32); err == nil {
		return uint(intValue)
	}
	return uint(defaultValue)
}

func (resource *armTemplateResource) getTagValue(key deployment.LookupTagKey, defaultValue string) string {
	if resource.Tags != nil {
		if value, exists := resource.Tags[string(key)]; exists {
			return value
		}
	}
	return defaultValue
}
