using System;

namespace Modm.Deployments
{
	public record CreateDeploymentResponse
	{
		public int Id { get; set; }
		public string Status { get; set; }
	}
}

