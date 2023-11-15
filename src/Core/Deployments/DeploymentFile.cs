using System.Text.Json;
using Microsoft.Extensions.Configuration;
using Microsoft.Extensions.Logging;
using Modm.Extensions;

namespace Modm.Deployments
{
    public class DeploymentFile
	{
        public const string FileName = "deployment.json";

        public DateTimeOffset Timestamp
        {
            get
            {
                return new FileInfo(GetDeploymentFilePath()).LastWriteTimeUtc;
            }
        }

        private static readonly JsonSerializerOptions serializerOptions = new JsonSerializerOptions
        {
            PropertyNamingPolicy = JsonNamingPolicy.CamelCase,
            WriteIndented = true
        };

        private readonly IConfiguration configuration;
        private readonly ILogger<DeploymentFile> logger;

        public DeploymentFile(IConfiguration configuration, ILogger<DeploymentFile> logger)
		{
            this.configuration = configuration;
            this.logger = logger;
        }

        public async Task<Deployment> Read(CancellationToken cancellationToken = default)
        {
            var path = GetDeploymentFilePath();

            if (!File.Exists(path))
            {
                this.logger.LogError($"{path} was not found");
                return new Deployment
                {
                    Id = 0,
                    Timestamp = DateTimeOffset.UtcNow,
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
            this.logger.LogInformation($"Writing Deployment json - {json}");
            await File.WriteAllTextAsync(GetDeploymentFilePath(), json, cancellationToken);
            this.logger.LogInformation($"Wrote the deployment json to {GetDeploymentFilePath()}");
        }

        private string GetDeploymentFilePath()
        {
            return Path.GetFullPath(Path.Combine(configuration.GetHomeDirectory(), FileName));
        }
    }
}

