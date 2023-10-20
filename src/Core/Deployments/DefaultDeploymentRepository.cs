using System;
using System.Threading;
using JenkinsNET.Models;
using Microsoft.Extensions.Logging;
using Modm.Engine.Jenkins.Client;

namespace Modm.Deployments
{
    public class DefaultDeploymentRepository : IDeploymentRepository
	{
        private IJenkinsClient client;
        private readonly JenkinsClientFactory clientFactory;
        private readonly DeploymentFile file;
        private readonly ILogger<DefaultDeploymentRepository> logger;

        public DefaultDeploymentRepository(JenkinsClientFactory clientFactory, DeploymentFile file, ILogger<DefaultDeploymentRepository> logger)
		{
            this.clientFactory = clientFactory;
            this.client = clientFactory.Create().GetAwaiter().GetResult();
            this.file = file;
            this.logger = logger;
        }

        protected IJenkinsClient JenkinsClient
        {
            get
            {
                if (this.client == null)
                {
                    this.client = this.clientFactory.Create().GetAwaiter().GetResult();
                }
                return this.client;
            }
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

