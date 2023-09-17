using System;
using MediatR;

namespace Modm.Engine.Notifications
{
	public class DeploymentStarted : INotification
	{
		public int Id { get; set; }
	}
}

