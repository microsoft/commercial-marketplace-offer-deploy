package mapper

import (
	"fmt"
	"log"

	"github.com/microsoft/commercial-marketplace-offer-deploy/internal/data"
	"github.com/microsoft/commercial-marketplace-offer-deploy/pkg/api"
	"github.com/mitchellh/mapstructure"
)

type CreateDeploymentMapper struct {
}

func NewCreateDeploymentMapper() *CreateDeploymentMapper {
	return &CreateDeploymentMapper{}
}

func (m *CreateDeploymentMapper) Map(from *api.CreateDeployment) (*data.Deployment, error) {
	err := m.validate(from)
	if err != nil {
		return nil, err
	}

	template := from.Template.(map[string]any)
	stages := m.getStages(template)

	deployment := &data.Deployment{
		Name:           *from.Name,
		Status:         "New",
		SubscriptionId: *from.SubscriptionID,
		ResourceGroup:  *from.ResourceGroup,
		Location:       *from.Location,
		Template:       template,
		Stages:         stages,
	}
	return deployment, nil
}

// create a method that validates the input
func (m *CreateDeploymentMapper) validate(from *api.CreateDeployment) error {
	_, ok := from.Template.(map[string]any)
	if !ok {
		err := fmt.Errorf("invalid .Template field value. Could not cast to map[string]any for '%s'", *from.Name)
		log.Printf("Error: %v", err)
		return err
	}
	return nil
}

// using the map as a graph, drill into each resource that's of type deployment
// then extract the values from the tags if they exist, otherwise use defaults to set the stage's fields
func (m *CreateDeploymentMapper) getStages(template map[string]any) []data.Stage {
	armTemplate := &armTemplate{}
	mapstructure.Decode(template, &armTemplate)

	stages := []data.Stage{}

	for _, resource := range armTemplate.Resources {
		if resource.isDeploymentResourceType() {
			stage := data.Stage{
				BaseWithGuidPrimaryKey: data.BaseWithGuidPrimaryKey{ID: resource.getId()},
				Name:                   resource.getName(),
				Status:                 "New",
			}
			stages = append(stages, stage)
		}
	}
	return stages
}
