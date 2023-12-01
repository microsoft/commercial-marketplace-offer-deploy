using System;
using Azure.ResourceManager;
using ClientApp.Notifications;
using MediatR;
using Microsoft.Extensions.Configuration;
using Microsoft.Extensions.Hosting;
using Microsoft.Extensions.Logging;
using Microsoft.Extensions.Options;
using Modm.Azure.Notifications;
using Modm.Deployments;

namespace ClientApp.Backend
{
	public class AzureDeploymentCleanupService : BackgroundService
    {
        private readonly DeploymentClient deploymentClient;
        private DateTime autoDeleteTime;
        private readonly IMediator mediator;
        private readonly IConfiguration configuration;
        private bool executeSelfDelete;

        private const string InstalledTimeKey = "InstalledTime";
        private const string ExpireInKey = "ExpireIn";

        public AzureDeploymentCleanupService(
            DeploymentClient deploymentClient,
            IMediator mediator,
            IConfiguration configuration)
		{
            this.deploymentClient = deploymentClient;
            this.configuration = configuration;
            this.autoDeleteTime = CalculateAutoDeleteTime();
        }

        private DateTime CalculateAutoDeleteTime()
        {
            var installedTimeString = this.configuration[InstalledTimeKey];
            if (!DateTime.TryParse(installedTimeString, out var installedTime))
            {
                installedTime = DateTime.UtcNow;
            }

            var expireInString = this.configuration[ExpireInKey];
            if (int.TryParse(expireInString, out var expireInMinutes) && expireInMinutes > -1)
            {
                this.executeSelfDelete = true;
                return installedTime.AddMinutes(expireInMinutes);
            }

            return installedTime.AddHours(24); 
        }


        protected override async Task ExecuteAsync(CancellationToken stoppingToken)
        {
            while (executeSelfDelete && !stoppingToken.IsCancellationRequested)
            {
                if (DateTime.UtcNow > autoDeleteTime)
                {
                    var deploymentsResponse = await this.deploymentClient.GetDeploymentInfo();
                    var resourceGroup = deploymentsResponse.Deployment.ResourceGroup;

                    var deleteInitiated = new DeleteInitiated(resourceGroup);
                    await mediator.Send(deleteInitiated);
                    break;
                }
                await Task.Delay(TimeSpan.FromMinutes(1), stoppingToken);
            }
        }
    }
}

