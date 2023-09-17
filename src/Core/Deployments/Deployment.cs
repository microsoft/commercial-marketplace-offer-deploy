using System;
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
	}
}

