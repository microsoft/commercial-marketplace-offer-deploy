using System;
using MediatR;
using Microsoft.Extensions.Hosting;
using Microsoft.Extensions.Logging;
using Modm.Deployments;
using Modm.Engine.Notifications;

namespace Modm.Engine
{
	public class DeploymentMonitorService : BackgroundService
	{
        private readonly IDeploymentEngine engine;
        private readonly ILogger<DeploymentMonitorService> logger;
        private bool deploymentStarted;
        private int id;

        public DeploymentMonitorService(IDeploymentEngine engine, ILogger<DeploymentMonitorService> logger)
        {
            this.engine = engine;
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
            var deployment = await engine.Get();

            // wait for the deployment to complete
            while (deployment.Status != DeploymentStatus.Completed)
            {
                await Task.Delay(1000, cancellationToken);
            }

            // deployment is complete
            logger.LogInformation("Deployment [{id}] completed at: {time}", id, DateTimeOffset.Now);
        }

        class DeploymentStartedHandler : INotificationHandler<DeploymentStarted>
        {
            private readonly DeploymentMonitorService service;

            public DeploymentStartedHandler(DeploymentMonitorService service)
            {
                this.service = service;
            }

            public Task Handle(DeploymentStarted notification, CancellationToken cancellationToken)
            {
                service.deploymentStarted = true;
                service.id = notification.Id;

                return Task.CompletedTask;
            }
        }
    }
}

