using System;
namespace Modm.ServiceHost
{
	public class ArtifactsWatcherOptions
	{
		public required string ArtifactsPath { get; set; }
		public required string DeploymentsUrl { get; set; }
	}
}

