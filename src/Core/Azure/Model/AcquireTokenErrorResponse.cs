// Copyright (c) Microsoft Corporation.
// Licensed under the MIT license.
﻿using System;
using System.Text.Json.Serialization;

namespace Modm.Azure
{
	record AcquireTokenErrorResponse
    {
        public const string InvalidRequestError = "invalid_request";

        [JsonPropertyName("error")]
        public required string Error { get; set; }

        [JsonPropertyName("error_description")]
        public required string ErrorDescription { get; set; }
    }
}

