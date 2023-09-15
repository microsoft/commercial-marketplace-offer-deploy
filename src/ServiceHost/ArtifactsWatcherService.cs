using System;
using System.Text;
using Newtonsoft.Json;
using MediatR;
using Modm.Deployments;
using Modm.ServiceHost.Notifications;
using Modm.Azure.Model;
using Modm.Azure;

namespace Modm.ServiceHost
{
    public class ArtifactsWatcherService : BackgroundService
    {
        const int DefaultWaitDelaySeconds = 10;
        private readonly IMetadataService metadataService;
        private ILogger<ArtifactsWatcherService> logger;

        ArtifactsWatcherOptions? options;
        bool controllerStarted;

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

            if (options == null)
            {
                throw new InvalidOperationException("Cannot start artifacts watcher. Options are null");
            }

            string base64UserData = "";

            while (true)
            {
                var instanceData = await this.metadataService.GetAsync();
                base64UserData = instanceData.Compute.UserData;
                if (!string.IsNullOrEmpty(base64UserData))
                {
                    break;
                }
                
                await Task.Delay(1000);
            }

            byte[] data = Convert.FromBase64String(base64UserData);
            string jsonString = Encoding.UTF8.GetString(data);
            UserData userData = JsonConvert.DeserializeObject<UserData>(jsonString);
            if (userData == null)
            {
                throw new InvalidDataException("The userData on the virtual machine instance is null");
            }

            if (userData.IsValid())
            {
                var response = await StartDeployment(userData.ArtifactsUri);
            }
        }

        private async Task<CreateDeploymentResponse> StartDeployment(string uri)
        {
            var request = new CreateDeploymentRequest { ArtifactsUri = uri };
            HttpResponseMessage response = await this.httpClient.PostAsJsonAsync(this.options?.DeploymentsUrl, request);
            response.EnsureSuccessStatusCode();

            this.logger.LogInformation("HTTP Post to [{url}] successful.", this.options?.DeploymentsUrl);

            return await response.Content.ReadAsAsync<CreateDeploymentResponse>();
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

