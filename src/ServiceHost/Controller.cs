using Ductus.FluentDocker.Model.Common;
using Ductus.FluentDocker.Builders;
using Ductus.FluentDocker.Services;
using Ductus.FluentDocker.Model.Compose;
using Modm.Configuration;
using Modm.Azure;
using MediatR;
using Modm.ServiceHost.Notifications;
using Modm.Extensions;

namespace Modm.ServiceHost
{
    /// <summary>
    /// The "Controller" that starts, monitors, and gracefully terminates modm
    /// in a background child process
    /// </summary>
    class Controller
    {
        #region Fields and Properties

        private readonly ControllerOptions options;
        private readonly IConfiguration configuration;
        private readonly IMediator mediator;
        private readonly ILogger<Controller> logger;
        ICompositeService? composeService;
        readonly IManagedIdentityService managedIdentityService;
        private readonly IMetadataService metadataService;
        private readonly IHostEnvironment environment;

        string EnvFilePath
        {
            get
            {
                return Path.Combine(options.ComposeFileDirectory, ".env");
            }
        }

        ICompositeService ComposeService
        {
            get
            {
                if (composeService == null)
                {
                    throw new NullReferenceException("Compose Service is null.");
                }
                return composeService;
            }
        }

        #endregion

        public Controller(ControllerOptions options, 
            IManagedIdentityService managedServiceIdentity,
            IMetadataService metadataService,
            IHostEnvironment environment, 
            IConfiguration configuration, 
            IMediator mediator, 
            ILogger<Controller> logger)
        {
            this.options = options;
            this.configuration = configuration;
            this.mediator = mediator;
            this.logger = logger;
            this.managedIdentityService = managedServiceIdentity;
            this.metadataService = metadataService;
            this.environment = environment;
        }

        public async Task StartAsync(CancellationToken cancellationToken = default)
        {
            logger.LogInformation("FQDN: {fqdn}", options.Fqdn);

            await UpdateEnvFileAsync();
            await StartComposeAsync();
            await Notify();

            logger.LogInformation("Running at: {time}", DateTimeOffset.Now);
            logger.LogInformation("Docker Compose state: {state}", ComposeService.State);
        }

        public Task StopAsync(CancellationToken cancellationToken = default)
        {
            ComposeService.Stop();
            return Task.CompletedTask;
        }

        /// <summary>
        /// Notify that the controller has started
        /// </summary>
        /// <returns></returns>
        async Task Notify()
        {
            var port = ComposeService.Containers.First(c => c.Name == "modm")
                .GetConfiguration()
                .NetworkSettings.Ports
                .First(p => p.Value != null && p.Value.FirstOrDefault() != null)
                .Value.First().Port;

            await mediator.Publish(new ControllerStarted
            {
                DeploymentsUrl = $"http://localhost:{port}/api/deployments",
                PackagePath = configuration.GetHomeDirectory(),
                StateFilePath = options.StateFilePath
            });
        }

        private async Task UpdateEnvFileAsync()
        {
            var userData = await metadataService.GetUserData();

            using var envFile = await GetEnvFileAsync();

            envFile.Set("MODM_HOME", configuration.GetHomeDirectory());

            var password = environment.IsDevelopment() ? "admin" : options.MachineName;
            envFile.Set("DEFAULT_ADMIN_PASSWORD", password);

            // app config endpoint for web host to read from
            envFile.Set("Azure__AppConfigEndpoint", userData.AppConfigEndpoint);

            // set for caddy to work
            envFile.Set("SITE_ADDRESS", options.Fqdn);
            envFile.Set("ACME_ACCOUNT_EMAIL", "nowhere@nowhere.com");

            var info = await managedIdentityService.GetAsync();

            // required by container environments to have managed identity flow from vm to container
            envFile.Set("AZURE_CLIENT_ID", info.ClientId.ToString());
            envFile.Set("AZURE_TENANT_ID", info.TenantId.ToString());
            envFile.Set("AZURE_SUBSCRIPTION_ID", info.SubscriptionId.ToString());

            await envFile.SaveAsync();
        }

        async Task StartComposeAsync()
        {
            await BuildComposeServiceAsync();
            ComposeService.Start();
        }

        private async Task BuildComposeServiceAsync()
        {
            var builder = new Builder()
                        .UseContainer()
                        .UseCompose()
                        .AssumeComposeVersion(ComposeVersion.V2)
                        .FromFile((TemplateString)options.ComposeFilePath)
                        .ServiceName("modm-service");

            using var envFile = await GetEnvFileAsync();
            logger.LogInformation("Environment Variables found: {count}", envFile.Items.Count);

            if (await envFile.AnyAsync())
            {
                logger.LogInformation("Adding variables from env file.");
                builder.WithEnvironment(envFile.ToArray());
            }

            var service = builder.RemoveOrphans().Build();

            service.StateChange += (object sender, StateChangeEventArgs e) =>
            {
                logger.LogInformation("Docker Compose state changed: {state}", e.State);
            };

            this.composeService = service;
        }

        private async Task<EnvFile> GetEnvFileAsync()
        {
            var envFile = EnvFile.New(EnvFilePath);
            await envFile.ReadAsync();

            return envFile;
        }

    }
}