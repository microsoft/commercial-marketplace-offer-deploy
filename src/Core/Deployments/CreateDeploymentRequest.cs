using System;

namespace Modm.Deployments
{
	public record CreateDeploymentRequest
	{
		/// <summary>
		/// The location of where the artifacts to be used for deployment/install, e.g. content.zip file that was in the app.zip
		/// </summary>
		public required string ArtifactsUri { get; set; }

		/// <summary>
		/// The deployment parameters
		/// </summary>
		public IDictionary<string,object> Parameters { get; set; }
	}
}

