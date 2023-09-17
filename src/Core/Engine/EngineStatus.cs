using System;
using Modm.Deployments;

namespace Modm.Engine
{
	public class EngineStatus
	{
		/// <summary>
		/// whether the engine is currently running a deployment
		/// </summary>
		public string IsRunning { get; set; }

		/// <summary>
		/// Info about the current deployment (if the engine is running this will be set)
		/// </summary>
		public Deployment Deployment { get; set; }	
	}
}

