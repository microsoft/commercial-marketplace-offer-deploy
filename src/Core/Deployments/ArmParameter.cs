using System;
using System.Collections.Generic;
using System.Text.Json;
using System.Text.Json.Serialization;

namespace Modm.Deployments
{
    public static class ArmParameterExtensions
    {
        public static string ToArmParametersJson(this IDictionary<string, object> inputParameters)
        {
            var armParameters = new ArmParameters
            {
                Parameters = new Dictionary<string, ArmParameter>()
            };

            var options = new JsonSerializerOptions { WriteIndented = true };

            foreach (var param in inputParameters)
            {
                var jsonValue = JsonSerializer.SerializeToElement(param.Value, options);
                armParameters.Parameters.Add(param.Key, new ArmParameter { Value = jsonValue });
            }

            return JsonSerializer.Serialize(armParameters, options);
        }
    }

    public class ArmParameter
	{
		[JsonPropertyName("value")]
		public JsonElement Value { get; set; }
	}

	public class ArmParameters
	{
        [JsonPropertyName("parameters")]
        public Dictionary<string, ArmParameter> Parameters { get; set; }
    }
}

