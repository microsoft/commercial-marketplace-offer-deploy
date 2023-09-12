using Ductus.FluentDocker.Model.Common;
using Ductus.FluentDocker.Builders;
using Ductus.FluentDocker.Services;
using Ductus.FluentDocker.Model.Compose;
using Modm.Configuration;

namespace Modm.ServiceHost
{
    /// <summary>
    /// The "Controller" that starts, monitors, and gracefully terminates modm
    /// in a background child process
    /// </summary>
    class Controller
    {
        private readonly ControllerOptions options;
        private readonly ILogger<Controller> logger;
        ICompositeService? composeService;

        public Controller(ControllerOptions options)
        {
            if (options.Logger == null)
            {
                throw new ArgumentNullException(nameof(options), "Logger cannot be null.");
            }

            this.options = options;
            this.logger = options.Logger;
        }

        public async Task StartAsync(CancellationToken cancellationToken = default)
        {
            logger.LogInformation("FQDN: {fqdn}", options.Fqdn);

            this.composeService = BuildDockerComposeService();
            composeService.StateChange += Service_StateChange;
            composeService.Start();

            options.Watcher?.Start();

            while (!cancellationToken.IsCancellationRequested)
            {
                logger.LogInformation("Running at: {time}", DateTimeOffset.Now);
                logger.LogInformation("Docker Compose state: {state}", composeService.State);

                await Task.Delay(10000, cancellationToken);
            }
        }

        public Task StopAsync(CancellationToken cancellationToken = default)
        {
            composeService?.Stop();
            return Task.CompletedTask;
        }

        private void Service_StateChange(object sender, StateChangeEventArgs e)
        {
            logger.LogInformation("Docker Compose state changed: {state}", e.State);
        }

        private ICompositeService BuildDockerComposeService()
        {
            var builder = new Builder()
                        .UseContainer()
                        .UseCompose()
                        .AssumeComposeVersion(ComposeVersion.V2)
                        .FromFile((TemplateString)options.ComposeFilePath);

            var envFilePath = Path.Combine(options.ComposeFileDirectory, ".env");
            var isEnvFileNextToComposeFile = File.Exists(envFilePath);

            if (isEnvFileNextToComposeFile)
            {
                var envFile = EnvFile.FromPath(envFilePath);
                if (envFile.HasItems)
                {
                    builder.WithEnvironment(envFile.Items.Select(item => $"{item.Key}={item.Value}").ToArray());
                }
            }

            var compositeService = builder.RemoveOrphans()
                        .WaitForHttp("jenkins", "http://localhost:8080/login")
                        .Build();

            return compositeService;
        }

    }
}