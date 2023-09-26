using System;
using MediatR;

namespace Modm.Deployments
{
    public record CreateDeploymentDefinition : StartDeploymentRequest, IRequest<DeploymentDefinition>
	{
        internal CreateDeploymentDefinition(StartDeploymentRequest request)
        {
            this.ArtifactsUri = request.ArtifactsUri;
            this.ArtifactsSig = request.ArtifactsSig;
            this.Parameters = request.Parameters;
        }
    }
}

