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
            var isBuilding = await GetBuild(cancellationToken);

            // wait for the deployment to complete
            while (isBuilding)
            {
                try
                {
                    isBuilding = await GetBuild(cancellationToken);
                }
                catch (Exception ex)
                {
                    logger.LogError(ex, "Exception thrown while trying to get build status");
                }
                await Task.Delay(1000, cancellationToken);
            }

            // deployment is complete
            logger.LogInformation("Deployment [{id}] completed at: {time}", id, DateTimeOffset.Now);
            Reset();
        }

        void Reset()
        {
            deploymentStarted = false;
            id = 0;
            name = string.Empty;
        }

        async Task<bool> GetBuild(CancellationToken cancellationToken)
        {
            try
            {
                var build = await client.Builds.GetAsync<JenkinsBuildBase>(name, id.ToString(), cancellationToken);
                if (build == null)
                {
                    return false;
                }

                return build.Building.GetValueOrDefault(false);
            }
            catch (Exception ex)
            {
                logger.LogError(ex, "Exception thrown while trying to get build status");
            }

            return false;
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

