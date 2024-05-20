using System;
namespace Modm.Deployments
{
	public record StartRedeploymentResult
	{
		public Deployment Deployment { get; set; }
		public List<string> Errors { get; set; }
	}
}

