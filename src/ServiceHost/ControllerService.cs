using Modm.Azure;
using Modm.ServiceHost.Extensions;

namespace Modm.ServiceHost
{
    public class ControllerService : BackgroundService
    {
        Controller? controller;
        private readonly IMetadataService metadataService;
        private readonly ILogger<ControllerService> logger;
        private readonly IServiceProvider serviceProvider;
        private readonly IConfiguration config;

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

            controller = ControllerBuilder.Create(this.logger)
                .UseFqdn(await metadataService.GetFqdnAsync())
                .UseComposeFile(GetComposeFilePath())
                .UsingServiceProvider(serviceProvider)
                .Build();

            await controller.StartAsync(stoppingToken);
        }

        public override Task StopAsync(CancellationToken cancellationToken)
        {
            controller?.StopAsync(cancellationToken);
            return base.StopAsync(cancellationToken);
        }

        string GetComposeFilePath()
        {
            return Path.Combine(config.GetHomeDirectory(), "service/docker-compose.yml");
        }
    }
}

