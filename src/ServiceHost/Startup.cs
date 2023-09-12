using System;
using System.Threading;
using Modm.Azure;

namespace Modm.ServiceHost
{
    public class Startup : BackgroundService
    {
        Controller? controller;
        private readonly InstanceMetadataService metadataService;
        private readonly ILogger<Startup> logger;
        private readonly IServiceProvider serviceProvider;
        private readonly ArtifactsWatcher watcher;

        public Startup(ArtifactsWatcher watcher, InstanceMetadataService metadataService, ILogger<Startup> logger, IServiceProvider serviceProvider)
        {
            this.watcher = watcher;
            this.metadataService = metadataService;
            this.logger = logger;
            this.serviceProvider = serviceProvider;
        }

        protected override async Task ExecuteAsync(CancellationToken stoppingToken)
        {
            logger.LogInformation("ServiceHost started.");

            controller = ControllerBuilder.Create(this.logger)
                .UseFqdn(await GetFqdnAsync())
                .UseComposeFile(GetComposeFilePath())
                .UseArtifactsWatcher(watcher)
                .UsingServiceProvider(serviceProvider)
                .Build();

            await controller.StartAsync(stoppingToken);
        }

        public override Task StopAsync(CancellationToken cancellationToken)
        {
            controller?.StopAsync(cancellationToken);
            return base.StopAsync(cancellationToken);
        }

        private async Task<string> GetFqdnAsync()
        {
            var isProduction = Environment.GetEnvironmentVariable("DOTNET_ENVIRONMENT") == "production";

            if (isProduction)
            {
                var metadata = await metadataService.GetAsync();
                return metadata.Fqdn;
            }
            return "localhost";
        }

        string GetComposeFilePath()
        {
            var folderPath = Environment.GetEnvironmentVariable("MODM_HOME") ?? throw new InvalidOperationException("MODM_HOME is null.");
            return Path.Combine(folderPath, "docker-compose.yml");
        }
    }
}

