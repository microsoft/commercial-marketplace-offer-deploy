package mapper

import (
	"github.com/google/uuid"
	"github.com/microsoft/commercial-marketplace-offer-deploy/pkg/deployment"
)

const deploymentResourceTypeName = "Microsoft.Resources/deployments"

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

func (resource *armTemplateResource) getTagValue(key deployment.LookupTagKey, defaultValue string) string {
	if resource.Tags != nil {
		if value, exists := resource.Tags[string(key)]; exists {
			return value
		}
	}
	return defaultValue
}
