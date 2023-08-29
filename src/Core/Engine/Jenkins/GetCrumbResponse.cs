using System;
using System.Text.Json.Serialization;

namespace Modm.Engine.Jenkins
{
    public record GetCrumbResponse
    {
        [JsonPropertyName("_class")]
        public required string Class { get; set; }

        [JsonPropertyName("crumb")]
        public required string Crumb { get; set; }

        [JsonPropertyName("crumbRequestField")]
        public required string RequestField { get; set; }
    }

}

