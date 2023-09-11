using System;
using System.Text.Json.Serialization;

namespace Modm.Engine
{
    [JsonConverter(typeof(JsonStringEnumConverter))]
    public enum EngineType
	{
		Jenkins,
		Bicep,
		Arm,
	}
}

