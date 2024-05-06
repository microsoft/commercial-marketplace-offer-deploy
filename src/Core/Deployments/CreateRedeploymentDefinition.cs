using System;
using MediatR;

namespace Modm.Deployments
{
    public record CreateRedeploymentDefinition : StartRedeploymentRequest, IRequest<DeploymentDefinition>
	{
        public new int DeploymentId { get; set; }

        internal CreateRedeploymentDefinition(int deploymentId, StartRedeploymentRequest request)
        {
            this.DeploymentId = deploymentId;
            this.Parameters = request.Parameters;
        }
    }
}

