package mapper

import (
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/to"
	"github.com/microsoft/commercial-marketplace-offer-deploy/internal/model"
	"github.com/microsoft/commercial-marketplace-offer-deploy/sdk"
)

type DeploymentMapper struct {
}

func NewDeploymentMapper() *DeploymentMapper {
	return &DeploymentMapper{}
}

func (m *DeploymentMapper) MapAll(deployments []model.Deployment) []sdk.Deployment {
	result := []sdk.Deployment{}

	for _, deployment := range deployments {
		result = append(result, m.Map(&deployment))
	}
	return result
}

func (m *DeploymentMapper) Map(deployment *model.Deployment) sdk.Deployment {
	result := sdk.Deployment{
		ID:             to.Ptr(int32(deployment.ID)),
		Name:           &deployment.Name,
		Location:       to.Ptr(deployment.Location),
		ResourceGroup:  to.Ptr(deployment.ResourceGroup),
		SubscriptionID: to.Ptr(deployment.SubscriptionId),
		Template:       to.Ptr(deployment.Template),
	}

	for _, stage := range deployment.Stages {
		result.Stages = append(result.Stages, &sdk.DeploymentStage{
			Name:           to.Ptr(stage.Name),
			ID:             to.Ptr(stage.ID.String()),
			Attributes:     stage.Attributes,
			Retries:        to.Ptr(int32(stage.Retries)),
			DeploymentName: to.Ptr(stage.AzureDeploymentName),
		})
	}
	return result
}
