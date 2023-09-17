using System;
using System.Text.Json.Serialization;

namespace Modm.Deployments
{
	/// <summary>
	/// The deployment
	/// </summary>
	public class Deployment
	{
		public int Id { get; set; }

		public string Status { get; set; }

		public DeploymentDefinition Definition { get; set; }

		public string Source { get; set; }

        // TODO: handle redeploy for failures.
        // For now, we are going to deploy based on either undefined (new) or if there was a failure
        public bool IsStartable => Status == DeploymentStatus.Undefined || Status == DeploymentStatus.Failure;
    }
}

