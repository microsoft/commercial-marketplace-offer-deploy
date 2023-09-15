using System;
using MediatR;

namespace Modm.ServiceHost.Notifications
{
	/// <summary>
	/// internal notification that the controller has been started
	/// </summary>
	public class ControllerStarted : INotification
	{
		public required string DeploymentsUrl { get; set; }
		public required string ArtifactsPath { get; set; }
	}
}

