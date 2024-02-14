// Copyright (c) Microsoft Corporation.
// Licensed under the MIT license.
﻿using System;
using System.Net.NetworkInformation;
using System.Text.Json.Serialization;

namespace Modm.Engine
{
	/// <summary>
	/// Information about the engine
	/// </summary>
	public record EngineInfo
	{
		[JsonPropertyName("isHealthy")]
        public required bool IsHealthy { get; set; }

        [JsonPropertyName("message")]
        public required string Message { get; set; }

        [JsonPropertyName("version")]
		public required string Version { get; set; }

		public static EngineInfo Default()
        {
			return new EngineInfo { IsHealthy = false, Version = "Unknown", Message = string.Empty };
        }
	}
}