using System;
using System.Text.Json.Serialization;

namespace Modm.Azure
{
    class TokenResponse
    {
        [JsonPropertyName("access_token")]
        public required string AccessToken { get; set; }

        [JsonPropertyName("client_id")]
        public Guid ClientId { get; set; }

        [JsonPropertyName("expires_in")]
        public long ExpiresIn { get; set; }

        [JsonPropertyName("expires_on")]
        public long ExpiresOn { get; set; }

        [JsonPropertyName("ext_expires_in")]
        public long ExtExpiresIn { get; set; }

        [JsonPropertyName("not_before")]
        public long NotBefore { get; set; }

        [JsonPropertyName("resource")]
        public required Uri Resource { get; set; }

        [JsonPropertyName("token_type")]
        public required string TokenType { get; set; }
    }

}

