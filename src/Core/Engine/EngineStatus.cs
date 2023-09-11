using System;
namespace Modm.Engine
{
	public class EngineStatus
	{
		public required bool IsHealthy { get; set; }
		public required EngineType EngineType { get; set; }
		public required string Version { get; set; }
	}
}

