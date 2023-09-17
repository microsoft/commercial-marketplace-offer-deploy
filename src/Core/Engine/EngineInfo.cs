using System;
using System.Net.NetworkInformation;

namespace Modm.Engine
{
	/// <summary>
	/// Information about the engine
	/// </summary>
	public record EngineInfo
	{
		public required bool IsHealthy { get; set; }

		public required EngineType EngineType { get; set; }

		public required string Version { get; set; }

		public static EngineInfo Default()
        {
			return new EngineInfo { EngineType = EngineType.Jenkins, IsHealthy = false, Version = "Unknown" };
        }
	}
}