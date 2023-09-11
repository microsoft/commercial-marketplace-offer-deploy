using System;
using Modm.Azure;
using Modm.Networking;

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

            var controller = ControllerBuilder.Create()
                .UseFqdn(await GetFqdn())
                .Build();

            await controller.StartAsync(stoppingToken);
        }

        private async Task<string> GetFqdn()
        {
            var metadata = await metadataService.GetAsync();
            return FqdnFactory.Create(metadata.Compute.ResourceId, metadata.Compute.Location);
        }
    }
}

