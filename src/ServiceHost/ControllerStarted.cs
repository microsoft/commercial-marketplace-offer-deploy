using System;
using MediatR;

namespace Modm.ServiceHost
{
	public class ControllerStarted : INotification
	{
		public required string DeploymentsUrl { get; set; }
		public required string ArtifactsPath { get; set; }
	}
}

