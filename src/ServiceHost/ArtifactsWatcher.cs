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
        private readonly string artifactsFilePath;
        private readonly string statusEndpoint;
        private readonly FileSystemWatcher fileWatcher;

        public ArtifactsWatcher(string artifactsFilePath, string statusEndpoint)
		{
            this.artifactsFilePath = string.IsNullOrEmpty(artifactsFilePath) ? GetDefaultArtifactsPath() : artifactsFilePath;
            this.statusEndpoint = statusEndpoint;

            fileWatcher = new FileSystemWatcher(Path.GetDirectoryName(artifactsFilePath));
            fileWatcher.Filter = Path.GetFileName(artifactsFilePath);
            fileWatcher.Created += OnFileCreated;
        }

        private string GetDefaultArtifactsPath()
        {
            string? modmHome = Environment.GetEnvironmentVariable("MODM_HOME");
            if (string.IsNullOrEmpty(modmHome))
            {
                throw new InvalidOperationException("$MODM_HOME environment variable is not set.");
            }
            return Path.Combine(modmHome, "artifacts.uri");
        }

        public void Start()
        {
            this.fileWatcher.EnableRaisingEvents = true;
        }

        private async void OnFileCreated(object sender, FileSystemEventArgs e)
        {
            try
            {
                string uri = File.ReadAllText(artifactsFilePath);

                if (Uri.IsWellFormedUriString(uri, UriKind.Absolute))
                {
                    bool isReady = await WaitForServiceReady();

                    if (isReady)
                    {
                        await SendHttpPost(uri);
                        Console.WriteLine("HTTP Post sent successfully.");
                    }
                    else
                    {
                        Console.WriteLine("External service is not ready.");
                    }
                }
                else
                {
                    Console.WriteLine("Invalid URI format.");
                }
            }
            catch (Exception ex)
            {
                Console.WriteLine($"Error: {ex.Message}");
            }
        }

        private async Task<bool> WaitForServiceReady()
        {
            int maxAttempts = 12; // 12 attempts * 5 seconds each = 1 minute
            int attemptIntervalSeconds = 5;

            for (int attempt = 0; attempt < maxAttempts; attempt++)
            {
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
            using (HttpClient client = new HttpClient())
            {
                HttpResponseMessage response = await client.GetAsync(statusEndpoint);
                if (response.IsSuccessStatusCode)
                {
                    // Assuming the response contains JSON representing EngineStatus
                    EngineStatus status = await response.Content.ReadAsAsync<EngineStatus>();
                    return status.IsHealthy;
                }

                return false;
            }
        }

        private async Task<CreateDeploymentResponse> SendHttpPost(string uri)
        {
            var request = new CreateDeploymentRequest { ArtifactsUri = uri };
            using (HttpClient client = new HttpClient())
            {
                HttpResponseMessage response = await client.PostAsJsonAsync(uri, request);
                response.EnsureSuccessStatusCode();
                return await response.Content.ReadAsAsync<CreateDeploymentResponse>();
            }
        }
    }
}

