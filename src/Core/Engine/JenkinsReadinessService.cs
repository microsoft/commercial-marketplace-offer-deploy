using System;
using Microsoft.Extensions.Hosting;
using Microsoft.Extensions.Logging;
using Microsoft.Extensions.Options;
using Modm.Jenkins;
using Modm.Jenkins.Client;
using Polly;
using Polly.Retry;

namespace Modm.Engine
{
    /// <summary>
    /// Monitors Jenkins to determine when it is ready to serve API requests
    /// </summary>
	public class JenkinsReadinessService : BackgroundService
    {
        private readonly string baseJenkinsUrl = "http://localhost:8080";
        private const int DefaultWaitDelaySeconds = 30;
        private const int MillisecondsInASecond = 1000;

        private readonly IDeploymentEngine engine;
        private readonly HttpClient httpClient;
        private readonly ILogger<JenkinsReadinessService> logger;

        private readonly AsyncRetryPolicy asyncRetryPolicy;

        private readonly JenkinsOptions jenkinsOptions;
        private bool isHealthy;
        private EngineInfo engineInfo;

        public JenkinsReadinessService(
            IDeploymentEngine engine,
            HttpClient httpClient,
            IOptions<JenkinsOptions> options,
            ILogger<JenkinsReadinessService> logger)
		{
            this.engine = engine;
            this.httpClient = httpClient;
            this.jenkinsOptions = options.Value;
            this.logger = logger;
            this.isHealthy = false;
            this.engineInfo = EngineInfo.Default();

            this.asyncRetryPolicy = Policy
               .Handle<Exception>()
               .WaitAndRetryForeverAsync(retryAttempt => TimeSpan.FromSeconds(Math.Pow(2, retryAttempt)));
        }

        protected override async Task ExecuteAsync(CancellationToken stoppingToken)
        {
            while (!stoppingToken.IsCancellationRequested)
            {
                if (await IsJenkinsLoginAvailableAsync())
                {
                    var engineInfo = await GetEngineInfoAsync();
                    if (engineInfo != null && engineInfo.IsHealthy != this.isHealthy)
                    {
                        this.isHealthy = engineInfo.IsHealthy;
                        PublishEngineInfo(engineInfo);
                    }
                }
                else
                {
                    PublishEngineInfo(EngineInfo.Default());
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

        private void PublishEngineInfo(EngineInfo engineInfo)
        {
            this.engineInfo = engineInfo;
        }

        private async Task<EngineInfo> GetEngineInfoAsync()
        {
            return await this.engine.GetInfo();
        }

        public EngineInfo GetEngineInfo()
        {
            return this.engineInfo;
        }
    }
}

