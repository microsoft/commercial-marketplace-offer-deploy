using System;
using MediatR;

namespace Modm.Deployments
{
    public record CreateDeploymentDefinition : StartDeploymentRequest, IRequest<DeploymentDefinition>
	{
        internal CreateDeploymentDefinition(StartDeploymentRequest request)
        {
            this.PackageUri = request.PackageUri;
            this.PackageHash = request.PackageHash;
            this.Parameters = request.Parameters;
        }
    }
}

