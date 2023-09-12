using System;
namespace Modm.ServiceHost
{
    public record CreateDeploymentRequest
    {
        public required string ArtifactsUri { get; set; }
    }
}

