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
        const int DefaultWaitDelaySeconds = 30;
        const int MaxAttempts = 10;
        private const int MillisecondsInASecond = 1000;

        private readonly IMetadataService metadataService;
        private readonly ILogger<PackageWatcherService> logger;

        private UserData? userData;

        PackageWatcherOptions? options;

        bool controllerStarted;
        int attempts = 0;

        private readonly HttpClient httpClient;

        public PackageWatcherService(IMetadataService metadataService, HttpClient httpClient, ILogger<PackageWatcherService> logger)
		{
            this.metadataService = metadataService;
            this.httpClient = httpClient;
            this.logger = logger;
            this.userData = null;
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
            if (options == null)
            {
                throw new InvalidOperationException("Cannot start installer package watcher. Options are null");
            }

            
            if (this.userData == null)
            {
                var base64UserData = await FetchBase64UserData(cancellation);
                if (!String.IsNullOrEmpty(base64UserData))
                {
                    this.userData = UserData.Deserialize(base64UserData) ?? throw new InvalidDataException("The userData on the virtual machine instance is null");
                }
                else
                {
                    logger.LogWarning("Unable to start deployment. Attempt {attempt}", attempts);
                    attempts++;
                    await Task.Delay(DefaultWaitDelaySeconds * MillisecondsInASecond, cancellation);
                    return false;
                }
            }

            return await TryStartDeployment(cancellation);
        }

        private async Task<string?> FetchBase64UserData(CancellationToken cancellation)
        {
            while (!cancellation.IsCancellationRequested)
            {
                var instanceData = await this.metadataService.GetAsync();
                if (!string.IsNullOrEmpty(instanceData.Compute.UserData))
                {
                    return instanceData.Compute.UserData;
                }

                await Task.Delay(DefaultWaitDelaySeconds * MillisecondsInASecond, cancellation);
            }

            return null;
        }

        private async Task<bool> TryStartDeployment(CancellationToken cancellation)
        {
            try
            {
                if (this.userData == null || !this.userData.IsValid())
                {
                    logger.LogError("Invalid UserData.");
                    return false;
                }

                logger.LogInformation("UserData was valid");

                if (this.options == null || String.IsNullOrEmpty(this.options.StateFilePath))
                {
                    logger.LogError("No StateFilePath present");
                    return false;
                }

                logger.LogInformation($"stateFilePath - {this.options.StateFilePath}");
                if (File.Exists(this.options.StateFilePath))
                {
                    return true;
                }

                var request = new StartDeploymentRequest
                {
                    PackageUri = userData.InstallerPackage.Uri,
                    PackageHash = userData.InstallerPackage.Hash,
                    Parameters = userData.Parameters ?? new Dictionary<string, object>()
                };

                var engineChecker = new EngineChecker("http://localhost:5000", this.httpClient);
                var isHealthy = await engineChecker.IsEngineHealthy();
                if (!isHealthy)
                {
                    return false;
                }

                await SubmitDeployment(request, cancellation);

                // Serialize the request to JSON and write it to the state file
                var json = JsonSerializer.Serialize(request, new JsonSerializerOptions { WriteIndented = true });
                await File.WriteAllTextAsync(this.options.StateFilePath, json);
                return true;
            }
            catch (Exception ex)
            {
                logger.LogError(ex, "Error deserializing UserData or starting deployment.");
                return false;
            }
        }

        private async Task SubmitDeployment(StartDeploymentRequest request, CancellationToken cancellation)
        {
            while (true)
            {
                try
                {
                    var response = await StartDeployment(request);
                    logger.LogInformation("Received deployment result, Id: {id}", response?.Deployment.Id);
                    return;
                }
                catch (Exception ex)
                {
                    logger.LogError(ex, "Error submitting deployment.");
                    await Task.Delay(DefaultWaitDelaySeconds * MillisecondsInASecond, cancellation);
                }
            }
        }

        private async Task<StartDeploymentResult?> StartDeployment(StartDeploymentRequest request)
        {
            HttpResponseMessage response = await this.httpClient.PostAsJsonAsync(this.options?.DeploymentsUrl, request);
            response.EnsureSuccessStatusCode();

            this.logger.LogInformation("HTTP Post to [{url}] successful.", this.options?.DeploymentsUrl);

            return await System.Text.Json.JsonSerializer.DeserializeAsync<StartDeploymentResult>(
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
                    DeploymentsUrl = notification.DeploymentsUrl,
                    StateFilePath = notification.StateFilePath
                };
                service.controllerStarted = true;

                return Task.CompletedTask;
            }

        }
    }
}

