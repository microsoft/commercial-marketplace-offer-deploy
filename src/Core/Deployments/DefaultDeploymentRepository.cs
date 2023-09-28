using System;
using System.Threading;
using JenkinsNET.Models;
using Microsoft.Extensions.Logging;
using Modm.Engine.Jenkins.Client;

namespace Modm.Deployments
{
    public class DefaultDeploymentRepository : IDeploymentRepository
	{
        private readonly IJenkinsClient client;
        private readonly DeploymentFile file;
        private readonly ILogger<DefaultDeploymentRepository> logger;

        public DefaultDeploymentRepository(IJenkinsClient client, DeploymentFile file, ILogger<DefaultDeploymentRepository> logger)
		{
            this.client = client;
            this.file = file;
            this.logger = logger;
        }

        public async Task<Deployment> Get(CancellationToken cancellationToken = default)
        {
            var deployment = await file.Read(cancellationToken);
            deployment.Status = await client.GetBuildStatus(deployment.Definition.DeploymentType, deployment.Id);
            deployment.IsStartable = await client.IsJobRunningOrWasAlreadyQueued(deployment.Definition.DeploymentType);

            return deployment;
        }

    }
}

