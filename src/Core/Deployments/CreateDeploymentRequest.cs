using System;

namespace Modm.Deployments
{
	public record CreateDeploymentRequest
	{
		public required string ArtifactsUri { get; set; }
	}
}

