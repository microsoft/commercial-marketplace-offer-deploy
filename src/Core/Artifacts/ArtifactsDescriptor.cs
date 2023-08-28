using System;
using Modm.Deployments;

namespace Modm.Artifacts
{
	public class ArtifactsDescriptor
	{
		public required DeploymentDefinition Definition { get; set; }

		/// <summary>
		/// the location path of the extracted artifacts which came from app.zip
		/// </summary>
        public required string FolderPath { get; set; }
	}
}

