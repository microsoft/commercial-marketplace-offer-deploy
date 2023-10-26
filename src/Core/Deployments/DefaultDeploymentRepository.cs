using Microsoft.Extensions.Logging;
using Modm.Jenkins.Client;

namespace Modm.Deployments
{
    public class DefaultDeploymentRepository : IDeploymentRepository
	{
        private readonly JenkinsClientFactory clientFactory;
        private readonly DeploymentFile file;
        private readonly ILogger<DefaultDeploymentRepository> logger;

        public DefaultDeploymentRepository(JenkinsClientFactory clientFactory, DeploymentFile file, ILogger<DefaultDeploymentRepository> logger)
		{
            this.clientFactory = clientFactory;
            this.file = file;
            this.logger = logger;
        }

        public async Task<Deployment> Get(CancellationToken cancellationToken = default)
        {
            logger.LogTrace("Fetching deployment information");

            using var client = await clientFactory.Create();

            var deployment = await file.Read(cancellationToken);
            deployment.Status = await client.GetBuildStatus(deployment.Definition.DeploymentType, deployment.Id);
            deployment.IsStartable = await client.IsJobRunningOrWasAlreadyQueued(deployment.Definition.DeploymentType);

            return deployment;
        }
    }
}

