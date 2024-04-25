using System;
using MediatR;

namespace Modm.Deployments
{
    public record CreateRedeploymentDefinition : StartDeploymentRequest, IRequest<DeploymentDefinition>
	{
        public int DeploymentId { get; set; }

        internal CreateRedeploymentDefinition(int deploymentId, StartDeploymentRequest request)
        {
            this.DeploymentId = deploymentId;
            this.PackageUri = request.PackageUri;
            this.PackageHash = request.PackageHash;
            this.Parameters = request.Parameters;
        }
    }
}

