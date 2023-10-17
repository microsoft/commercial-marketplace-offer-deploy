using System;
using System.Text.Json.Serialization;
using MediatR;
using Modm.Packaging;
using Modm.Serialization;

namespace Modm.Deployments
{
	public record StartDeploymentRequest : IRequest<StartDeploymentResult>, IRequest<DeploymentDefinition>
    {
		/// <summary>
		/// The location of the installer package to be used for deployment/install, e.g. installer.zip file that was in the app.zip
		/// </summary>
		public string PackageUri { get; set; }

		/// <summary>
		/// The origional signature of the installer package, used to verify it hasn't been tampered with
		/// </summary>
		public string PackageHash { get; set; }

		/// <summary>
		/// The deployment parameters
		/// </summary>
		[JsonConverter(typeof(DictionaryStringObjectJsonConverter))]
		public Dictionary<string,object> Parameters { get; set; }


        /// <summary>
        /// Gets the installer package uri as an <see cref="Packaging.PackageUri"/>
        /// </summary>
        /// <returns></returns>
        public PackageUri GetUri()
		{
			return new PackageUri(PackageUri);
        }
	}
}

