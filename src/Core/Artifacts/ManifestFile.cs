using System;
using System.Text.Json;
using Modm.Deployments;

namespace Modm.Artifacts
{
	public class ManifestFile
    {
        /// <summary>
        /// the name of the manifest file.
        /// </summary>
        /// <remarks>
        /// This value must never change as it's a convention that the Partner Center CLI must conform to
        /// when generating the artifacts file
        /// </remarks>
        public const string FileName = "manifest.json";

        /// <summary>
        /// 
        /// </summary>
        /// <param name="directoryPath">The directory where the manifest file is located</param>
        /// <returns></returns>
        /// <exception cref="FileNotFoundException"></exception>
        public static async Task<DeploymentDefinition> Read(string directoryPath)
		{
            var filePath = Path.Combine(directoryPath, FileName);

            if (!File.Exists(filePath))
            {
                throw new FileNotFoundException($"{FileName} not found in '{filePath}'.");
            }

            var json = await File.ReadAllTextAsync(filePath);
            var definition = JsonSerializer.Deserialize<DeploymentDefinition>(json, new JsonSerializerOptions
            {
                PropertyNamingPolicy = JsonNamingPolicy.CamelCase
            });

            return definition;
        }
	}
}

