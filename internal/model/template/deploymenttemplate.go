package template

import "github.com/microsoft/commercial-marketplace-offer-deploy/pkg/deployment"

type DeploymentTemplate struct {
	source     map[string]any
	lookupTags map[string][]deployment.LookupTag
}

// used to capture nested template elements
type nestedTemplateElement map[string]any

func NewDeploymentTemplate(source map[string]any) *DeploymentTemplate {
	return &DeploymentTemplate{
		source:     source,
		lookupTags: map[string][]deployment.LookupTag{},
	}
}

// tags a nested template with a lookup tag
func (t *DeploymentTemplate) Tag(name string, lookupTag deployment.LookupTag) {
	t.lookupTags[name] = append(t.lookupTags[name], lookupTag)
}

// gets the template
func (t *DeploymentTemplate) Build() map[string]interface{} {
	if t.source == nil {
		return nil
	}

	template := t.source

	nestedTemplates := t.getNestedTemplates()
	if len(nestedTemplates) == 0 {
		return template
	}

	for _, nestedTemplate := range nestedTemplates {
		name := nestedTemplate["name"].(string)
		if lookupTags, ok := t.lookupTags[name]; ok {
			t.addLookupTags(nestedTemplate, lookupTags)
		}
	}

	return template //return the modified template
}

// adds lookup tags to a nested template
func (t *DeploymentTemplate) addLookupTags(element nestedTemplateElement, lookupTags []deployment.LookupTag) {
	if tagsEntry, ok := element["tags"]; ok {
		if tagsMap, ok := tagsEntry.(map[string]any); ok {
			for _, lookupTag := range lookupTags {
				tagsMap[string(lookupTag.Key)] = lookupTag.Value
			}
		}
	}
}

// gets all nested template elements at the root of the template's "resources" section
func (t *DeploymentTemplate) getNestedTemplates() []nestedTemplateElement {
	nestedTemplates := []nestedTemplateElement{}
	template := t.source

	// set the operationId on every nested template at level 1 of the parent template
	if resourcesEntry, ok := template["resources"]; ok {
		if resources, ok := resourcesEntry.([]any); ok {
			if len(resources) > 0 {
				for _, resourceEntry := range resources {
					if resourceMap, ok := resourceEntry.(map[string]any); ok {
						if isDeploymentResourceType(resourceMap) {
							nestedTemplates = append(nestedTemplates, resourceMap)
						}
					}
				}
			}
		}
	}
	return nestedTemplates
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
