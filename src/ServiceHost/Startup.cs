using System;
using Modm.Azure;

namespace Modm.ServiceHost
{
    public class Startup : BackgroundService
    {
        private readonly InstanceMetadataService metadataService;
        private readonly ILogger<Worker> _logger;

        public Startup(InstanceMetadataService metadataService, ILogger<Worker> logger)
        {
            this.metadataService = metadataService;
            _logger = logger;
        }

        protected override async Task ExecuteAsync(CancellationToken stoppingToken)
        {
            while (!stoppingToken.IsCancellationRequested)
            {
                var fqdn = await metadataService.GetFqdnAsync();
                _logger.LogInformation("FQDN: {fqdn}", fqdn);

                _logger.LogInformation("Running at: {time}", DateTimeOffset.Now);
                await Task.Delay(10000, stoppingToken);
            }
        }
    }
}

