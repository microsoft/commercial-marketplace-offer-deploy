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
            var fqdn = await metadataService.GetFqdnAsync();
            _logger.LogInformation("FQDN: {fqdn}", fqdn);

            while (!stoppingToken.IsCancellationRequested)
            {
                _logger.LogInformation("Worker running at: {time}", DateTimeOffset.Now);
                await Task.Delay(1000, stoppingToken);
            }
        }
    }
}

