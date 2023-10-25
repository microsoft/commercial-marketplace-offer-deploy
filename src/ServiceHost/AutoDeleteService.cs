using System;
using System.Net.Http;
using System.IO;
using System.Threading;
using System.Threading.Tasks;
using SharpCompress.Common;
using Ductus.FluentDocker.Commands;
using Modm.Azure;
using Modm.Configuration;
using Modm.Extensions;

namespace Modm.ServiceHost
{
	public class AutoDeleteService : BackgroundService
    {

        private readonly IMetadataService metadataService;

        private readonly HttpClient httpClient;

        private readonly IConfiguration configuration;

        // The date and time when the action should be triggered.
        public DateTime AutoDeleteTime { get; set; }

        private readonly ILogger<AutoDeleteService> logger;

        public AutoDeleteService(
            IMetadataService metadataService,
            IConfiguration configuration,
            HttpClient httpClient,
            ILogger<AutoDeleteService> logger)
		{
            this.metadataService = metadataService;
            this.configuration = configuration;
            this.httpClient = httpClient;
            this.logger = logger;

            InitializeAutoDeleteTime();
        }

        private void InitializeAutoDeleteTime()
        {
            string stateFilePath = Path.Combine(this.configuration.GetHomeDirectory(), "autodelete.txt");

            if (!File.Exists(stateFilePath))
            {
                File.Create(stateFilePath);
            }

            var fileCreationDate = File.GetCreationTime(stateFilePath);
            // Set the ActionTime to be 24 hours from the read date.
            AutoDeleteTime = fileCreationDate.AddHours(24);
        }

        protected override async Task ExecuteAsync(CancellationToken stoppingToken)
        {
            while (!stoppingToken.IsCancellationRequested)
            {
                // If the current time has passed the ActionTime
                if (DateTime.Now >= AutoDeleteTime)
                {
                    await AutoDelete();

                    // After calling AutoDelete, stop the service or 
                    // set a new ActionTime for the next action.
                    // For now, I'll just stop the loop to avoid multiple calls to AutoDelete.
                    break;
                }

                // Wait for a while before checking again.
                // You can adjust the delay as needed.
                await Task.Delay(TimeSpan.FromMinutes(1), stoppingToken);
            }
        }

        
        private async Task AutoDelete()
        {
            try
            {
                var instanceData = await this.metadataService.GetAsync();
                var resourceGroupName = instanceData.Compute.ResourceGroupName;
                var apiUrl = $"http://localhost:5000/{resourceGroupName}/deletemodmresources";

                // Send a POST request
                var response = await this.httpClient.PostAsync(apiUrl, null);  // null indicates no content in the body of the POST request.

                if (response.IsSuccessStatusCode)
                {
                    string responseBody = await response.Content.ReadAsStringAsync();
                    this.logger.LogInformation($"API Response: {responseBody}");
                }
                else
                {
                    this.logger.LogError($"Error calling the API: {response.StatusCode}");
                }
            }
            catch (Exception ex)
            {
                this.logger.LogError($"An error occurred: {ex.Message}");
            }

        }
    }
}

