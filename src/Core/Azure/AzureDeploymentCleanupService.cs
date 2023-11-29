using System;
using MediatR;
using Microsoft.Extensions.Hosting;
using Microsoft.Extensions.Logging;
using Modm.Azure.Notifications;

namespace Modm.Azure
{
	public class AzureDeploymentCleanupService : BackgroundService
    {
        private readonly IMediator mediator;
        private readonly string resourceGroupName;

        public AzureDeploymentCleanupService(IMediator mediator, string resourceGroupName)
		{
            this.mediator = mediator;
            this.resourceGroupName = resourceGroupName;
        }

        protected override async Task ExecuteAsync(CancellationToken stoppingToken)
        {
            var cleanupLimitReached = new CleanupLimitReached(this.resourceGroupName);
            await mediator.Send(cleanupLimitReached);
        }
    }
}

