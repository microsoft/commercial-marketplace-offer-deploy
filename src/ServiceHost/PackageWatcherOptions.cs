using System;
namespace Modm.ServiceHost
{
	public class PackageWatcherOptions
	{
		public required string PackagePath { get; set; }
		public required string DeploymentsUrl { get; set; }
	}
}

