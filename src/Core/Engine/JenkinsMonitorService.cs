using JenkinsNET.Models;
using MediatR;
using Microsoft.Extensions.Hosting;
using Microsoft.Extensions.Logging;
using Modm.Engine.Jenkins.Client;
using Modm.Engine.Notifications;

namespace Modm.Engine
{
    /// <summary>
    /// Monitors the active jenkins build triggered through the engine
    /// </summary>
	public class JenkinsMonitorService : BackgroundService
	{
        private readonly IJenkinsClient client;
        private readonly ILogger<JenkinsMonitorService> logger;

        private bool deploymentStarted;
        private int id;
        private string name;

        public JenkinsMonitorService(IJenkinsClient client, ILogger<JenkinsMonitorService> logger)
        {
            this.client = client;
            this.logger = logger;
        }

        protected override async Task ExecuteAsync(CancellationToken stoppingToken)
        {
            logger.LogInformation("Deployment monitor service started at: {time}", DateTimeOffset.Now);

            while (!stoppingToken.IsCancellationRequested)
            {
                await WaitUntilDeploymentHasStarted(stoppingToken);
                await MonitorDeployment(stoppingToken);
            }
        }

        async Task WaitUntilDeploymentHasStarted(CancellationToken cancellationToken)
        {
            logger.LogInformation("Waiting for deployment to start");

            while (!deploymentStarted)
            {
                await Task.Delay(1000, cancellationToken);
            }

            logger.LogInformation("Deployment [{id}] started at: {time}", id, DateTimeOffset.Now);
        }

        async Task MonitorDeployment(CancellationToken cancellationToken)
        {
            var build = await GetBuild(cancellationToken);

            // wait for the deployment to complete
            while (build.Building.GetValueOrDefault(false))
            {
                build = await GetBuild(cancellationToken);
                await Task.Delay(1000, cancellationToken);
            }

            // deployment is complete
            logger.LogInformation("Deployment [{id}] completed at: {time}", id, DateTimeOffset.FromUnixTimeSeconds(build.TimeStamp.Value));
            Reset();
        }

        void Reset()
        {
            deploymentStarted = false;
            id = 0;
            name = string.Empty;
        }

        async Task<JenkinsBuildBase> GetBuild(CancellationToken cancellationToken)
        {
            return await client.Builds.GetAsync<JenkinsBuildBase>(name, id.ToString(), cancellationToken);
        }

        public class DeploymentStartedHandler : INotificationHandler<DeploymentStarted>
        {
            private readonly JenkinsMonitorService service;

            public DeploymentStartedHandler(JenkinsMonitorService service)
            {
                this.service = service;
            }

            public Task Handle(DeploymentStarted notification, CancellationToken cancellationToken)
            {
                service.deploymentStarted = true;
                service.id = notification.Id;
                service.name = notification.Name;

                return Task.CompletedTask;
            }
        }
    }
}

