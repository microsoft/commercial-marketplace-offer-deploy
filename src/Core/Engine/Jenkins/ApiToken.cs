using System;
using System.Text.Json.Serialization;

namespace Modm.Engine.Jenkins
{
    internal record ApiTokenData
    {
        [JsonPropertyName("tokenName")]
        public required string Name { get; set; }

        [JsonPropertyName("tokenUuid")]
        public required string Uuid { get; set; }

        [JsonPropertyName("tokenValue")]
        public required string Value { get; set; }
    }

    /// <summary>
    /// The response payload that will be received from Jenkins when requesting
    /// to create a Jenkins API Token in order to begin making API calls to Jenkins
    /// </summary>
    internal record GenerateApiTokenResponse
    {
        [JsonPropertyName("status")]
        public required string Status { get; set; }

        [JsonPropertyName("data")]
        public required ApiTokenData Data { get; set; }
    }
}

