using System;
using System.Text.Json.Serialization;
using MediatR;
using Modm.Packaging;
using Modm.Serialization;

namespace Modm.Deployments
{
	public record StartRedeploymentRequest : IRequest<StartRedeploymentResult>, IRequest<DeploymentDefinition>
    {
		/// <summary>
		/// The ID of the deployment that you wish to redeploy
		/// </summary>
		public string DeploymentId { get; set; }

		/// <summary>
		/// The deployment parameters
		/// </summary>
		[JsonConverter(typeof(DictionaryStringObjectJsonConverter))]
		public Dictionary<string, object> Parameters { get; set; }

	}
}
