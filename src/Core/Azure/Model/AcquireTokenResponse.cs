// Copyright (c) Microsoft Corporation.
// Licensed under the MIT license.
﻿using System;
using System.Text.Json.Serialization;

namespace Modm.Azure
{
    record AcquireTokenResponse
    {
        [JsonPropertyName("access_token")]
        public required string AccessToken { get; set; }

        [JsonPropertyName("client_id")]
        public Guid ClientId { get; set; }

        [JsonPropertyName("expires_in")]
        public required string ExpiresIn { get; set; }

        [JsonPropertyName("expires_on")]
        public required string ExpiresOn { get; set; }

        [JsonPropertyName("ext_expires_in")]
        public required string ExtExpiresIn { get; set; }

        [JsonPropertyName("not_before")]
        public required string NotBefore { get; set; }

        [JsonPropertyName("resource")]
        public required Uri Resource { get; set; }

        [JsonPropertyName("token_type")]
        public required string TokenType { get; set; }
    }

}

