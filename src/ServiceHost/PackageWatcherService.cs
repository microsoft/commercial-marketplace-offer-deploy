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
        private const int MillisecondsInASecond = 1000;

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
            //if (attempts > MaxAttempts)
            //{
            //    logger.LogWarning("Max attempts reached while processing user data.");
            //    return true;
            //}

            if (options == null)
            {
                throw new InvalidOperationException("Cannot start installer package watcher. Options are null");
            }

            var base64UserData = await FetchBase64UserData(cancellation);

            if (string.IsNullOrEmpty(base64UserData))
            {
                logger.LogWarning("Unable to start deployment. Attempt {attempt}", attempts);
                attempts++;
                await Task.Delay(DefaultWaitDelaySeconds * MillisecondsInASecond, cancellation);
                return false;
            }

            return await TryStartDeployment(base64UserData, cancellation);
        }

        private async Task<string> FetchBase64UserData(CancellationToken cancellation)
        {
            //TODO: Determine if this is necessary or if we can count on the
            // metadata service being available
            while (true)
            {
                var instanceData = await this.metadataService.GetAsync();
                if (!string.IsNullOrEmpty(instanceData.Compute.UserData))
                {
                    return instanceData.Compute.UserData;
                }

                await Task.Delay(DefaultWaitDelaySeconds * MillisecondsInASecond, cancellation);
            }
        }

        private async Task<bool> TryStartDeployment(string base64UserData, CancellationToken cancellation)
        {
            try
            {
                var userData = UserData.Deserialize(base64UserData) ?? throw new InvalidDataException("The userData on the virtual machine instance is null");

                if (!userData.IsValid())
                {
                    logger.LogError("Invalid UserData.");
                    return false;
                }

                logger.LogInformation("UserData was valid");
                logger.LogInformation($"StateFilePath - {this.options?.StateFilePath}");

                if (File.Exists(this.options?.StateFilePath))
                {
                    return true;
                }
                else
                {
                    var request = new StartDeploymentRequest
                    {
                        PackageUri = userData.InstallerPackage.Uri,
                        PackageHash = userData.InstallerPackage.Hash,
                        Parameters = userData.Parameters ?? new Dictionary<string, object>()
                    };

                    // Serialize the request to JSON and write it to the state file
                    var json = JsonSerializer.Serialize(request, new JsonSerializerOptions { WriteIndented = true });
                    await File.WriteAllTextAsync(this.options?.StateFilePath, json);

                    await SubmitDeployment(request, cancellation);
                    return true;
                }
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
                    //var request = new StartDeploymentRequest
                    //{
                    //    PackageUri = userData.InstallerPackage.Uri,
                    //    PackageHash = userData.InstallerPackage.Hash,
                    //    Parameters = userData.Parameters ?? new Dictionary<string, object>()
                    //};

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

