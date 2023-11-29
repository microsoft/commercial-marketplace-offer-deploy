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
        private DateTime autoDeleteTime;

        public AzureDeploymentCleanupService(IMediator mediator, string resourceGroupName)
		{
            this.mediator = mediator;
            this.resourceGroupName = resourceGroupName;
            this.autoDeleteTime = GetDeployTime().AddHours(24);
        }

        private DateTime GetDeployTime()
        {
            return DateTime.UtcNow;
        }

        protected override async Task ExecuteAsync(CancellationToken stoppingToken)
        {
            while (!stoppingToken.IsCancellationRequested)
            {
                if (DateTime.UtcNow > autoDeleteTime)
                {
                    var cleanupLimitReached = new CleanupLimitReached(this.resourceGroupName);
                    await mediator.Send(cleanupLimitReached);
                    break;
                }
                await Task.Delay(TimeSpan.FromMinutes(1), stoppingToken);
            }
        }
    }
}

