using System;
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

        [JsonPropertyName("engineType")]
		public required EngineType EngineType { get; set; }

        [JsonPropertyName("version")]
		public required string Version { get; set; }

		public static EngineInfo Default()
        {
			return new EngineInfo { EngineType = EngineType.Jenkins, IsHealthy = false, Version = "Unknown" };
        }
	}
}