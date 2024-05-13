using System;
using System.Text.Json;
using Modm.Serialization;

namespace Modm.Deployments
{
    /// <summary>
    /// Writes out parameters for a deployment as tfvars file
    /// </summary>
	class TerraformParametersFile : IDeploymentParametersFile
	{
        public string FullPath
        {
            get {  return Path.Combine(destinationDirectory, FileName); }
        }

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
            if (Directory.Exists(destinationDirectory))
            {
                var json = JsonSerializer.Serialize(parameters, new JsonSerializerOptions
                {
                    Converters = { new DictionaryStringObjectJsonConverter() }
                });

                // Validate write permission and file creation
                using (FileStream fs = new FileStream(FullPath, FileMode.Create, FileAccess.Write))
                {
                    using (StreamWriter writer = new StreamWriter(fs))
                    {
                        await writer.WriteAsync(json);
                    }
                }
            }

        }

        public async Task Delete()
        {
            if (File.Exists(FullPath))
            {
                await Task.Run(() => File.Delete(FullPath));
            }
        }
    }
}

