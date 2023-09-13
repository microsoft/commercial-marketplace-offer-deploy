using System.Net.Http.Json;
using Modm.Deployments;

namespace Modm.ServiceHost
{
    public class ArtifactsWatcher
	{
        public const string ArtifactsUriFileName = "artifacts.uri";

        private readonly ILogger<ArtifactsWatcher> logger;
        private string? deploymentsUrl;
        private readonly FileSystemWatcher fileWatcher;
        private readonly HttpClient httpClient;

        public ArtifactsWatcher(HttpClient client, ILogger<ArtifactsWatcher> logger)
		{
            this.logger = logger;
            this.httpClient = client;

            fileWatcher = new FileSystemWatcher
            {
                Filter = ArtifactsUriFileName
            };
        }

        public Task StartAsync(ArtifactsWatcherOptions options)
        {
            logger.LogInformation("Starting artifacts watcher");

            if (!Uri.IsWellFormedUriString(options.DeploymentsUrl, UriKind.Absolute))
            {
                this.logger.LogError("Received invalid URI format as [{uri}]", options.DeploymentsUrl);
                throw new ArgumentException("options.DeploymentsUrl must be well formed uri.", nameof(options));
            }

            fileWatcher.Path = options.ArtifactsPath;
            fileWatcher.Created += OnFileCreated;

            deploymentsUrl = options.DeploymentsUrl;

            this.logger.LogInformation("Artifacts watcher started. Watching: {directory}", fileWatcher.Path);
            this.logger.LogInformation("Artifacts watcher filter: {filter}", fileWatcher.Filter);

            this.fileWatcher.EnableRaisingEvents = true;

            return Task.CompletedTask;
        }

        private async void OnFileCreated(object sender, FileSystemEventArgs e)
        {
            this.logger.LogInformation("File created event");
            try
            {;
                string uri = File.ReadAllText(e.FullPath);
                await StartDeployment(uri);
            }
            catch (Exception ex)
            {
                this.logger.LogError(ex, "Error handing file created");
            }
        }

        private async Task<CreateDeploymentResponse> StartDeployment(string uri)
        {
            var request = new CreateDeploymentRequest { ArtifactsUri = uri };
            HttpResponseMessage response = await this.httpClient.PostAsJsonAsync(this.deploymentsUrl, request);
            response.EnsureSuccessStatusCode();

            this.logger.LogInformation("HTTP Post to [{url}] successful.", this.deploymentsUrl);

            return await response.Content.ReadAsAsync<CreateDeploymentResponse>();
        }
    }
}

