using System;
using System.Text.Json.Serialization;
using MediatR;
using Modm.Artifacts;
using Modm.Serialization;

namespace Modm.Deployments
{
	public record StartDeploymentRequest : IRequest<StartDeploymentResult>, IRequest<DeploymentDefinition>
    {
		/// <summary>
		/// The location of where the artifacts to be used for deployment/install, e.g. content.zip file that was in the app.zip
		/// </summary>
		public string ArtifactsUri { get; set; }

		/// <summary>
		/// The deployment parameters
		/// </summary>
		[JsonConverter(typeof(DictionaryStringObjectJsonConverter))]
		public IDictionary<string,object> Parameters { get; set; }


		/// <summary>
		/// Gets the artifacts uri as an <see cref="ArtifactsUri"/>
		/// </summary>
		/// <returns></returns>
		public ArtifactsUri GetUri()
		{
			return new ArtifactsUri(ArtifactsUri);
        }
	}
}

