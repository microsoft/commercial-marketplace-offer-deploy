using Ductus.FluentDocker.Model.Common;
using Ductus.FluentDocker.Builders;
using Ductus.FluentDocker.Services;
using Ductus.FluentDocker.Model.Compose;
using Modm.Configuration;
using Modm.Azure;

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
        readonly ManagedIdentityService managedIdentityService;

        public Controller(ControllerOptions options)
        {
            if (options.Logger == null)
            {
                throw new ArgumentNullException(nameof(options), "Logger cannot be null.");
            }

            this.options = options;
            this.logger = options.Logger;
            this.managedIdentityService = options.ManagedIdentityService;
        }

        public async Task StartAsync(CancellationToken cancellationToken = default)
        {
            logger.LogInformation("FQDN: {fqdn}", options.Fqdn);

            await SetEnvFileAsync();

            StartCompose();
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

        void StartCompose()
        {
            this.composeService = BuildDockerComposeService();
            composeService.StateChange += (object sender, StateChangeEventArgs e) =>
            {
                logger.LogInformation("Docker Compose state changed: {state}", e.State);
            };
            composeService.Start();
        }

        private async Task SetEnvFileAsync()
        {
            var envFilePath = Path.Combine(options.ComposeFileDirectory, ".env");
            var envFile = EnvFileReader.FromPath(envFilePath);

            var writer = new EnvFileWriter(envFile.Items);

            // set for caddy
            writer.Add("SITE_ADDRESS", options.Fqdn);

            var info = await managedIdentityService.GetAsync();

            // required by container environments
            writer.Add("AZURE_CLIENT_ID", info.ClientId.ToString());
            writer.Add("AZURE_TENANT_ID", info.TenantId.ToString());
            writer.Add("AZURE_SUBSCRIPTION_ID", info.SubscriptionId.ToString());

            await writer.WriteAsync(envFilePath);
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
                var envFile = EnvFileReader.FromPath(envFilePath);
                if (envFile.HasItems)
                {
                    builder.WithEnvironment(envFile.Items.Select(item => $"{item.Key}={item.Value}").ToArray());
                }
            }

            var compositeService = builder.RemoveOrphans()
                        .WaitForHttp("jenkins", "http://localhost:8080/login", timeout: 60000, (response, attempt) =>
                        {
                            logger.LogInformation("Engine check [{attempt}]. HTTP Status [{statusCode}]", attempt, response.Code);
                            return response.Code == System.Net.HttpStatusCode.OK ? 0 : 500;
                        })
                        .Build();

            return compositeService;
        }

    }
}