using System;
using System;
using System.IO;
using System.Net;
using System.Net.Http;
using System.Net.Http.Json;
using System.Threading.Tasks;
using Modm.Engine;

namespace Modm.ServiceHost
{
	public class ArtifactsWatcher
	{
        public const string ArtifactsUriFileName = "artifacts.uri";

        private readonly ILogger<ArtifactsWatcher> logger;
        private readonly string statusEndpoint;
        private readonly FileSystemWatcher fileWatcher;
        private readonly HttpClient httpClient;

        public ArtifactsWatcher(HttpClient client, string artifactsFilePath, string statusEndpoint, ILogger<ArtifactsWatcher> logger)
		{
            this.statusEndpoint = statusEndpoint;
            this.logger = logger;
            this.httpClient = client;

            var expandedPath = Environment.ExpandEnvironmentVariables(artifactsFilePath);
            fileWatcher = new FileSystemWatcher(Path.GetDirectoryName(expandedPath));
            fileWatcher.Filter = ArtifactsUriFileName;
            fileWatcher.Created += OnFileCreated;
        }


        public void Start()
        {
            this.logger.LogInformation("Artifacts watcher started. Watching: {directory}", fileWatcher.Path);
            this.logger.LogInformation("Artifacts watcher filter: {filter}", fileWatcher.Filter);
            this.fileWatcher.EnableRaisingEvents = true;
        }

        private async void OnFileCreated(object sender, FileSystemEventArgs e)
        {
            this.logger.LogInformation("File created.");
            try
            {
                var expandedPath = Environment.ExpandEnvironmentVariables(e.FullPath);
                this.logger.LogInformation($"exmpanedPath: {expandedPath}");
                string uri = File.ReadAllText(expandedPath);
                this.logger.LogInformation("uri: {uri}", uri);
                if (Uri.IsWellFormedUriString(uri, UriKind.Absolute))
                {
                    await SendHttpPost(uri);
                    this.logger.LogInformation("HTTP Post sent successfully.");
                }
                else
                {
                    this.logger.LogInformation("Invalid URI format.");
                }
            }
            catch (Exception ex)
            {
                this.logger.LogError($"Error: {ex.Message}");
            }
        }

        private async Task<bool> WaitForServiceReady()
        {
            int maxAttempts = 12; // 12 attempts * 5 seconds each = 1 minute
            int attemptIntervalSeconds = 5;

            for (int attempt = 0; attempt < maxAttempts; attempt++)
            {
                this.logger.LogInformation("inside WaitForServiceReady loop.");
                bool isReady = await CheckServiceStatus();

                if (isReady)
                {
                    return true;
                }

                await Task.Delay(attemptIntervalSeconds * 1000);
            }

            throw new TimeoutException("Timed out waiting for the service to be ready.");
        }

        private async Task<bool> CheckServiceStatus()
        {
            this.logger.LogInformation("inside CheckServiceStatus.");
            HttpResponseMessage response = await this.httpClient.GetAsync(statusEndpoint);
            if (response.IsSuccessStatusCode)
            {
                EngineStatus status = await response.Content.ReadAsAsync<EngineStatus>();
                this.logger.LogInformation($"Engine status: {status}");
                return status.IsHealthy;
            }

            return false;
        }

        private async Task<CreateDeploymentResponse> SendHttpPost(string uri)
        {
            var request = new CreateDeploymentRequest { ArtifactsUri = uri };
            HttpResponseMessage response = await this.httpClient.PostAsJsonAsync(uri, request);
            response.EnsureSuccessStatusCode();
            return await response.Content.ReadAsAsync<CreateDeploymentResponse>();
        }
    }
}

