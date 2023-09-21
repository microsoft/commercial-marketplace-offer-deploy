using System.Text.Json;
using Microsoft.Extensions.Configuration;
using Modm.Extensions;

namespace Modm.Deployments
{
    public class DeploymentFile
	{
        public const string FileName = "deployment.json";

        private static readonly JsonSerializerOptions serializerOptions = new JsonSerializerOptions
        {
            PropertyNamingPolicy = JsonNamingPolicy.CamelCase,
            WriteIndented = true
        };

        private readonly IConfiguration configuration;

        public DeploymentFile(IConfiguration configuration)
		{
            this.configuration = configuration;
        }

        public async Task<Deployment> Read(CancellationToken cancellationToken = default)
        {
            var path = GetDeploymentFilePath();

            if (!File.Exists(path))
            {
                return new Deployment
                {
                    Id = 0,
                    Status = DeploymentStatus.Undefined
                };
            }

            var json = await File.ReadAllTextAsync(GetDeploymentFilePath(), cancellationToken);
            var deployment = JsonSerializer.Deserialize<Deployment>(json, serializerOptions);

            return deployment;
        }

        public async Task Write(Deployment deployment, CancellationToken cancellationToken)
        {
            var json = JsonSerializer.Serialize(deployment, serializerOptions);
            await File.WriteAllTextAsync(GetDeploymentFilePath(), json, cancellationToken);
        }

        private string GetDeploymentFilePath()
        {
            return Path.GetFullPath(Path.Combine(configuration.GetHomeDirectory(), FileName));
        }
    }
}

