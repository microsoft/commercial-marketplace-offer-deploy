using System;
using System.Text.Json;

namespace Modm.Deployments
{
    /// <summary>
    /// Writes out parameters for a deployment as tfvars file
    /// </summary>
	class TerraformParametersFile : IDeploymentParametersFile
	{
        /// <summary>
        /// This file name will automatically be picked up by terraform <see cref="https://developer.hashicorp.com/terraform/language/values/variables"/>
        /// </summary>
        public const string FileName = "terraform.tfvars.json";
        private readonly string destinationDirectory;

        public TerraformParametersFile(string destinationDirectory)
		{
            this.destinationDirectory = destinationDirectory;
        }

        public async Task Write(IDictionary<string, object> parameters)
        {
            var filePath = Path.Combine(destinationDirectory, FileName);
            var json = JsonSerializer.Serialize(parameters);

            await File.WriteAllTextAsync(filePath, json);
        }
    }
}

