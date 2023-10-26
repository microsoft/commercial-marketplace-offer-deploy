using System;
using System.Text.Json.Serialization;

namespace Modm.Jenkins.Model
{
    /// <summary>
    /// Created from JSON result from http://localhost:8080/crumbIssuer/api/json
    /// </summary>
    record GetCrumbResponse
    {
        [JsonPropertyName("_class")]
        public required string Class { get; set; }

        [JsonPropertyName("crumb")]
        public required string Crumb { get; set; }

        [JsonPropertyName("crumbRequestField")]
        public required string RequestField { get; set; }
    }

}

