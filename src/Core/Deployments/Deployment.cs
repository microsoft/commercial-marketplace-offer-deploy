using System;
using System.Text.Json.Serialization;

namespace Modm.Deployments
{
	/// <summary>
	/// The deployment instance
	/// </summary>
	public class Deployment
	{
		public int Id { get; set; }

		public string Status { get; set; }

		public string ResourceGroup { get; set; }

		public string SubscriptionId { get; set; }

		public DeploymentDefinition Definition { get; set; }

		public IEnumerable<DeploymentResource> Resources { get; set; }

		public bool IsStartable { get; internal set; }

		public Deployment()
		{
			Resources = new List<DeploymentResource>();
		}
	}
}

