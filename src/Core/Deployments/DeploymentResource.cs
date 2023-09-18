using System;
namespace Modm.Deployments
{
	public record DeploymentResource
	{
		public string Name { get; set; }
		public string Type { get; set; }
		public string State { get; set; }
		public DateTimeOffset Timestamp { get; set; }
	}
}

