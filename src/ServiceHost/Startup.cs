using System;
using Modm.Azure;

namespace Modm.ServiceHost
{
    public class Startup : BackgroundService
    {
        private readonly InstanceMetadataService metadataService;
        private readonly ILogger<Worker> logger;

        public Startup(InstanceMetadataService metadataService, ILogger<Worker> logger)
        {
            this.metadataService = metadataService;
            this.logger = logger;
        }

        protected override async Task ExecuteAsync(CancellationToken stoppingToken)
        {
            logger.LogInformation("ServiceHost started.");
            var metadata = await metadataService.GetAsync();

            var controller = ControllerBuilder.Create()
                .UseFqdn(metadata.Fqdn)
                .Build();

            await controller.StartAsync(stoppingToken);
        }
    }
}

