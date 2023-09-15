using Modm.Azure;
using MediatR;
using Modm.ServiceHost.Notifications;

namespace Modm.ServiceHost
{
    /// <summary>
    /// Monitors the availability of a managed identity being available. once it identifies
    /// it can acquire, it acquires the information and notifies internally within
    /// the service host
    /// </summary>
    public class ManagedIdentityMonitorService : BackgroundService
	{
        const int DefaultWaitDelaySeconds = 10;

        private readonly IManagedIdentityService managedIdentityService;
        private readonly IMediator mediator;

        public ManagedIdentityMonitorService(IManagedIdentityService managedIdentityService, IMediator mediator)
		{
            this.managedIdentityService = managedIdentityService;
            this.mediator = mediator;
        }

        protected override async Task ExecuteAsync(CancellationToken stoppingToken)
        {
            while (!await managedIdentityService.IsAccessibleAsync(stoppingToken))
            {
                await Task.Delay(DefaultWaitDelaySeconds * 1000, stoppingToken);
            }

            await mediator.Publish(new ManagedIdentityAcquired(), stoppingToken);
        }
    }
}

