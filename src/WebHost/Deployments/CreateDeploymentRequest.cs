using System;
namespace WebHost.Deployments
{
	public record CreateDeploymentRequest
	{
		public required string ArtifactsUri { get; set; }
	}
}

