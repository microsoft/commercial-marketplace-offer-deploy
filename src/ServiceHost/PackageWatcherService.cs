using MediatR;
using Modm.Deployments;
using Modm.ServiceHost.Notifications;
using Modm.Azure.Model;
using Modm.Azure;
using System.Text.Json;

namespace Modm.ServiceHost
{
    /// <summary>
    /// Watcher service that handles the installer package
    /// </summary>
    public class PackageWatcherService : BackgroundService
    {
        const int DefaultWaitDelaySeconds = 10;
        const int MaxAttempts = 10;

        private readonly IMetadataService metadataService;
        private readonly ILogger<PackageWatcherService> logger;

        PackageWatcherOptions? options;

        bool controllerStarted;
        int attempts = 0;

        private readonly HttpClient httpClient;

        public PackageWatcherService(IMetadataService metadataService, HttpClient httpClient, ILogger<PackageWatcherService> logger)
		{
            this.metadataService = metadataService;
            this.httpClient = httpClient;
            this.logger = logger;
        }

        protected override async Task ExecuteAsync(CancellationToken cancellationToken)
        {
            await WaitForControllerToStart(cancellationToken);

            var userDataProcessed = false;

            while (!userDataProcessed)
            {
                userDataProcessed = await TryToProcessUserData(cancellationToken);
            }
        }

        private async Task<bool> TryToProcessUserData(CancellationToken cancellation)
        {
            if (attempts > MaxAttempts)
            {
                logger.LogWarning("Max attempts reached while processing user data.");
                return true;
            }

            if (options == null)
            {
                throw new InvalidOperationException("Cannot start installer package watcher. Options are null");
            }

            string base64UserData;

            while (true)
            {
                var instanceData = await this.metadataService.GetAsync();
                base64UserData = instanceData.Compute.UserData;

                if (!string.IsNullOrEmpty(base64UserData))
                {
                    break;
                }

                await Task.Delay(DefaultWaitDelaySeconds * 1000, cancellation);
            }

            try
            {
                var userData = UserData.Deserialize(base64UserData) ?? throw new InvalidDataException("The userData on the virtual machine instance is null");

                if (userData.IsValid())
                {
                    logger.LogInformation("UserData was valid");

                    var request = new StartDeploymentRequest
                    {
                        PackageUri = userData.InstallerPackage.Uri,
                        PackageHash = userData.InstallerPackage.Hash,
                        Parameters = userData.Parameters ?? new Dictionary<string, object>()
                    };

                    
                    var response = await StartDeployment(request);
                    logger.LogInformation("Received deployment result, Id: {id}", response?.Deployment.Id);

                    return true;
                }
            }
            catch (Exception ex)
            {
                logger.LogError(ex, "Error deserializing UserData or starting deployment.");
            }

            logger.LogWarning("Unable to start deployment. Attempt {attempt}", attempts);

            attempts++;
            await Task.Delay(DefaultWaitDelaySeconds * 1000, cancellation);

            return false;
        }

        private async Task<StartDeploymentResult?> StartDeployment(StartDeploymentRequest request)
        {
            HttpResponseMessage response = await this.httpClient.PostAsJsonAsync(this.options?.DeploymentsUrl, request);
            response.EnsureSuccessStatusCode();

            this.logger.LogInformation("HTTP Post to [{url}] successful.", this.options?.DeploymentsUrl);

            return await JsonSerializer.DeserializeAsync<StartDeploymentResult>(
                response.Content.ReadAsStream(), new JsonSerializerOptions
                {
                    PropertyNameCaseInsensitive = true,
                    PropertyNamingPolicy = JsonNamingPolicy.CamelCase
                });
        }

        async Task WaitForControllerToStart(CancellationToken cancellationToken)
        {
            while (!controllerStarted)
            {
                await Task.Delay(DefaultWaitDelaySeconds * 1000, cancellationToken);
            }
        }

        class ControllerStartedHandler : INotificationHandler<ControllerStarted>
        {
            private readonly PackageWatcherService service;

            public ControllerStartedHandler(PackageWatcherService service)
            {
                this.service = service;
            }

            public Task Handle(ControllerStarted notification, CancellationToken cancellationToken)
            {
                service.options = new PackageWatcherOptions
                {
                    PackagePath = notification.PackagePath,
                    DeploymentsUrl = notification.DeploymentsUrl
                };
                service.controllerStarted = true;

                return Task.CompletedTask;
            }

        }
    }
}

