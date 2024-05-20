using System;
using System.Text.Json.Serialization;
using Modm.Serialization;

namespace ClientApp.Commands
{
	public class SubmitRedeployment
	{
        /// <summary>
        /// The ID of the Deployment you wish to redeploy
        /// </summary>
		public string DeploymentId { get; set; }

        /// <summary>
        /// The deployment parameters
        /// </summary>
        [JsonConverter(typeof(DictionaryStringObjectJsonConverter))]
        public Dictionary<string, object> Parameters { get; set; }
	}
}

