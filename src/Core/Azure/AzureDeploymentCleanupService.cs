using System;
using Azure.ResourceManager;
using Microsoft.Extensions.Hosting;
using Microsoft.Extensions.Logging;
using Microsoft.Extensions.Options;
using Modm.Azure.Notifications;
using Modm.Deployments;

namespace Modm.Azure
{
	public class AzureDeploymentCleanupService : BackgroundService
    {
        private readonly DeploymentClient deploymentClient;
        private readonly ArmClient armClient;
        private DateTime autoDeleteTime;

        public AzureDeploymentCleanupService(DeploymentClient deploymentClient, ArmClient armClient)
		{
            this.deploymentClient = deploymentClient;
            this.armClient = armClient;
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
                    var deploymentsResponse = await this.deploymentClient.GetDeploymentInfo();
                    var resourceGroup = deploymentsResponse.Deployment.ResourceGroup;
                    
                    var armCleanup = new AzureDeploymentCleanup(armClient);
                    bool deleted = await armCleanup.DeleteResourcePostDeployment(resourceGroup);

                    //var cleanupLimitReached = new CleanupLimitReached(resourceGroup);
                    //await mediator.Send(cleanupLimitReached);
                    break;
                }
                await Task.Delay(TimeSpan.FromMinutes(1), stoppingToken);
            }
        }
    }
}

