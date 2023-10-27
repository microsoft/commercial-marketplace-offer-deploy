using System;
using Microsoft.Extensions.Hosting;
using Microsoft.Extensions.Logging;
using Microsoft.Extensions.Options;
using Modm.Jenkins;
using Modm.Jenkins.Client;
using Polly;
using Polly.Utilities;
using Polly.Retry;

namespace Modm.Engine
{
    /// <summary>
    /// Monitors Jenkins to determine when it is ready to serve API requests
    /// </summary>
	public class JenkinsReadinessService : BackgroundService
    {
        private const int DefaultWaitDelaySeconds = 30;
        private const int MillisecondsInASecond = 1000;
        private const int MaxRetries = 6;

        private readonly IDeploymentEngine engine;
        private readonly HttpClient httpClient;
        private readonly ILogger<JenkinsReadinessService> logger;

        private readonly AsyncRetryPolicy asyncRetryPolicy;

        private readonly JenkinsOptions jenkinsOptions;
        private EngineInfo engineInfo;
        private readonly JenkinsClientFactory clientFactory;

        public JenkinsReadinessService(
            HttpClient httpClient,
            JenkinsClientFactory clientFactory,
            IOptions<JenkinsOptions> options,
            ILogger<JenkinsReadinessService> logger)
		{
            this.httpClient = httpClient;
            this.clientFactory = clientFactory;
            this.jenkinsOptions = options.Value;
            this.logger = logger;
            this.engineInfo = EngineInfo.Default();
            this.asyncRetryPolicy = Policy
               .Handle<Exception>()
               .WaitAndRetryAsync(MaxRetries, retryAttempt =>
                    TimeSpan.FromSeconds(Math.Pow(2, retryAttempt)));
        }

        protected override async Task ExecuteAsync(CancellationToken stoppingToken)
        {
            while (!stoppingToken.IsCancellationRequested)
            {
                if (await IsJenkinsLoginAvailableAsync())
                {
                    var engineInfo = await GetEngineInfoAsync();
                    UpdateEngineInfo(engineInfo);
                }
                else
                {
                    UpdateEngineInfo(EngineInfo.Default());
                }

                await Task.Delay(DefaultWaitDelaySeconds * MillisecondsInASecond, stoppingToken);
            }
        }

        private async Task<bool> IsJenkinsLoginAvailableAsync()
        {
            var jenkinsBaseUrl = this.jenkinsOptions.BaseUrl;
            var result = await this.asyncRetryPolicy.ExecuteAsync(async () =>
            {
                var response = await httpClient.GetAsync($"{jenkinsBaseUrl}/login");
                if (!response.IsSuccessStatusCode)
                {
                    this.logger.LogError($"The jenkins /login uri returned {response.StatusCode}");
                }
                return response.IsSuccessStatusCode;
            });

            return result;
        }

        private void UpdateEngineInfo(EngineInfo engineInfo)
        {
            this.engineInfo = engineInfo;
        }

        private async Task<EngineInfo> GetEngineInfoAsync()
        {
            var result = EngineInfo.Default();

            try
            {
                using var client = await this.clientFactory.Create();

                var info = await client.GetInfo();
                var node = await client.GetBuiltInNode();

                result.IsHealthy = !node.Offline;
                result.Message = $"Offline reason: {node.OfflineCauseReason}, Temporarily offline: {node.TemporarilyOffline}";

            }
            catch (Exception ex)
            {
                this.logger.LogError(ex, "error occurred fetching engine info.");
                result.Message = ex.Message;
            }

            return result;
        }

        public EngineInfo GetEngineInfo()
        {
            return this.engineInfo;
        }
    }
}

