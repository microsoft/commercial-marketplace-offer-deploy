using ClientApp.Notifications;
using MediatR;

namespace ClientApp.Backend
{
    public class AzureDeploymentCleanupService : BackgroundService
    {
        private DateTime autoDeleteTime;
        private readonly IMediator mediator;
        private readonly IConfiguration configuration;
        private bool executeSelfDelete;

        private const string InstalledTimeKey = "InstalledTime";
        private const string ExpireInKey = "ExpireIn";
        private const string ResourceGroupNameKey = "WEBSITE_RESOURCE_GROUP";

        public AzureDeploymentCleanupService(
            IMediator mediator,
            IConfiguration configuration)
		{
            this.configuration = configuration;
            this.mediator = mediator;
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
                    var resourceGroupName = configuration[ResourceGroupNameKey];
                    var deleteInitiated = new DeleteInitiated(resourceGroupName);

                    await mediator.Send(deleteInitiated);

                    break;
                }

                await Task.Delay(TimeSpan.FromMinutes(1), stoppingToken);
            }
        }
    }
}

