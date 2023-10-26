using MediatR;
using Modm.Azure;
using Modm.Extensions;
using Modm.ServiceHost.Notifications;

namespace Modm.ServiceHost
{
    public class ControllerService : BackgroundService
    {
        Controller? controller;
        private readonly IMetadataService metadataService;
        private readonly ILogger<ControllerService> logger;
        private readonly IServiceProvider serviceProvider;
        private readonly IConfiguration config;

        bool managedIdentityAcquired = false;

        public ControllerService(IMetadataService metadataService, IServiceProvider serviceProvider, IConfiguration config, ILogger<ControllerService> logger)
        {
            this.metadataService = metadataService;
            this.logger = logger;
            this.serviceProvider = serviceProvider;
            this.config = config;
        }

        protected override async Task ExecuteAsync(CancellationToken stoppingToken)
        {
            logger.LogInformation("ServiceHost started.");

            await WaitUntilManagedIdentityIsAcquired(stoppingToken);

            var metadata = await metadataService.GetAsync();

            controller = ControllerBuilder.Create(this.logger)
                .UseFqdn(await metadataService.GetFqdnAsync())
                .UseMachineName(metadata.Compute.Name)
                .UseComposeFile(GetComposeFilePath())
                .UseStateFile(GetStateFilePath())
                .UsingServiceProvider(serviceProvider)
                .Build();

            await controller.StartAsync(stoppingToken);
        }

        public override async Task StopAsync(CancellationToken cancellationToken)
        {
            if (controller != null)
                await controller.StopAsync(cancellationToken);

            await base.StopAsync(cancellationToken);
        }

        async Task WaitUntilManagedIdentityIsAcquired(CancellationToken cancellationToken)
        {
            logger.LogInformation("Waiting to initialize controller until managed identity is acquired");

            while (!managedIdentityAcquired)
            {
                await Task.Delay(1000, cancellationToken);
            }

            logger.LogInformation("managed identity acquired. proceeding with controller initialization");
        }

        string GetComposeFilePath()
        {
            return Path.Combine(config.GetHomeDirectory(), "service/docker-compose.yml");
        }

        string GetStateFilePath()
        {
            return Path.Combine(config.GetHomeDirectory(), "service/state.txt");
        }

        class ManagedIdentityAcquiredHandler : INotificationHandler<ManagedIdentityAcquired>
        {
            private readonly ControllerService service;

            public ManagedIdentityAcquiredHandler(ControllerService service)
            {
                this.service = service;
            }

            public Task Handle(ManagedIdentityAcquired notification, CancellationToken cancellationToken)
            {
                service.managedIdentityAcquired = true;
                return Task.CompletedTask;
            }
        }
    }
}

