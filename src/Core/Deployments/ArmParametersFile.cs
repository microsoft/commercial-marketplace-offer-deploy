using System.Text.Json;
using System.Text.Json.Serialization;
using Microsoft.WindowsAzure.ResourceStack.Common.Extensions;

namespace Modm.Deployments
{
    public class ArmParametersFile : IDeploymentParametersFile
    {
        public const string FileName = "parameters.json";
        private readonly string destinationDirectory;

        public string FullPath => Path.Combine(destinationDirectory, FileName);

        public ArmParametersFile(string destinationDirectory)
		{
            this.destinationDirectory = destinationDirectory;
		}

        public async Task Write(IDictionary<string, object> parameters)
        {
            var json = JsonSerializer.Serialize(new ArmParametersFileContent
            {
                Parameters = parameters?.ToDictionary(p => p.Key, p => ArmParameter.From(p.Value))
            });
            await File.WriteAllTextAsync(FullPath, json);
        }

        class ArmParametersFileContent
        {
            [JsonPropertyName("$schema")]
            public string Schema { get; } = "https://schema.management.azure.com/schemas/2019-04-01/deploymentParameters.json#";

            [JsonPropertyName("contentVersion")]
            public string ContentVersion { get; } = "1.0.0.0";

            [JsonPropertyName("parameters")]
            public Dictionary<string, ArmParameter> Parameters { get; set; }
        }

        class ArmParameter
        {
            static readonly JsonSerializerOptions options = new JsonSerializerOptions { WriteIndented = true };

            [JsonPropertyName("value")]
            public JsonElement Value { get; set; }

            public static ArmParameter From(object value)
            {
                return new ArmParameter { Value = JsonSerializer.SerializeToElement(value, options) };
            }
        }
    }
}

