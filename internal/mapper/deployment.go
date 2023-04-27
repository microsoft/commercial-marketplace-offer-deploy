package mapper

import (
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/to"
	"github.com/microsoft/commercial-marketplace-offer-deploy/internal/data"
	"github.com/microsoft/commercial-marketplace-offer-deploy/pkg/api"
)

type DeploymentMapper struct {
}

func NewDeploymentMapper() *DeploymentMapper {
	return &DeploymentMapper{}
}

func (m *DeploymentMapper) MapAll(deployments []data.Deployment) []api.Deployment {
	result := []api.Deployment{}

	for _, deployment := range deployments {
		result = append(result, m.Map(&deployment))
	}
	return result
}

func (m *DeploymentMapper) Map(deployment *data.Deployment) api.Deployment {
	result := api.Deployment{
		ID:     to.Ptr(int32(deployment.ID)),
		Name:   &deployment.Name,
		Status: &deployment.Status,
	}

	for _, stage := range deployment.Stages {
		result.Stages = append(result.Stages, &api.DeploymentStage{
			Name: to.Ptr(stage.Name),
			ID:   to.Ptr(stage.ID.String()),
		})
	}
	return result
}
