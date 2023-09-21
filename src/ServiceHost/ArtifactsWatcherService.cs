using System;
using System.Text;
using Newtonsoft.Json;
using MediatR;
using Modm.Deployments;
using Modm.ServiceHost.Notifications;
using Modm.Azure.Model;
using Modm.Azure;
using System.Security.Policy;
using System.Threading;

namespace Modm.ServiceHost
{
    public class ArtifactsWatcherService : BackgroundService
    {
        const int DefaultWaitDelaySeconds = 10;
        const int MaxAttempts = 10;

        private readonly IMetadataService metadataService;
        private ILogger<ArtifactsWatcherService> logger;

        ArtifactsWatcherOptions? options;
        bool controllerStarted;
        int attempts = 0;

        private readonly HttpClient httpClient;

        public ArtifactsWatcherService(IMetadataService metadataService, HttpClient httpClient, ILogger<ArtifactsWatcherService> logger)
		{
            this.metadataService = metadataService;
            this.httpClient = httpClient;
            this.logger = logger;
        }

        protected override async Task ExecuteAsync(CancellationToken cancellationToken)
        {
            await WaitForControllerToStart(cancellationToken);

            while (!cancellationToken.IsCancellationRequested)
            {
                await TryToProcessUserData(cancellationToken);
            }
        }

        private async Task TryToProcessUserData(CancellationToken cancellation)
        {
            if (attempts > MaxAttempts)
            {
                return;
            }

            if (options == null)
            {
                throw new InvalidOperationException("Cannot start artifacts watcher. Options are null");
            }

            string base64UserData;

            while (true)
            {
                var instanceData = await this.metadataService.GetAsync();
                base64UserData = instanceData?.Compute.UserData;
                if (!string.IsNullOrEmpty(base64UserData))
                {
                    break;
                }

                await Task.Delay(1000);
            }

            try
            {
                var userData = UserData.Deserialize(base64UserData) ?? throw new InvalidDataException("The userData on the virtual machine instance is null");

                if (userData.IsValid())
                {
                    var request = new StartDeploymentRequest
                    {
                        ArtifactsUri = userData.ArtifactsUri,
                        Parameters = userData.Parameters ?? new Dictionary<string, object>()
                    };

                    var response = await StartDeployment(request);
                }
            }
            catch (Exception ex)
            {
                logger.LogError(ex, "Error deserializing UserData or starting deployment.");
            }

            logger.LogWarning("Unable to start deployment. Attempt {attempt}", attempts);

            attempts++;
            await Task.Delay(DefaultWaitDelaySeconds * 1000, cancellation);

        }

        private async Task<StartDeploymentResult> StartDeployment(StartDeploymentRequest request)
        {
          
            HttpResponseMessage response = await this.httpClient.PostAsJsonAsync(this.options?.DeploymentsUrl, request);
            response.EnsureSuccessStatusCode();

            this.logger.LogInformation("HTTP Post to [{url}] successful.", this.options?.DeploymentsUrl);

            return await response.Content.ReadAsAsync<StartDeploymentResult>();
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
            private readonly ArtifactsWatcherService service;

            public ControllerStartedHandler(ArtifactsWatcherService service)
            {
                this.service = service;
            }

            public Task Handle(ControllerStarted notification, CancellationToken cancellationToken)
            {
                service.options = new ArtifactsWatcherOptions
                {
                    ArtifactsPath = notification.ArtifactsPath,
                    DeploymentsUrl = notification.DeploymentsUrl
                };
                service.controllerStarted = true;

                return Task.CompletedTask;
            }

        }
    }
}

