using System;
using System.Net.Http;
using System.IO;
using System.Threading;
using System.Threading.Tasks;
using SharpCompress.Common;
using Ductus.FluentDocker.Commands;
using Modm.Azure;

namespace Modm.ServiceHost
{
	public class AutoDeleteService : BackgroundService
    {
        private readonly string stateFilePath;

        private readonly IMetadataService metadataService;

        private readonly HttpClient httpClient;

        // The date and time when the action should be triggered.
        public DateTime AutoDeleteTime { get; set; }

        private readonly ILogger<AutoDeleteService> logger;

        public AutoDeleteService(string stateFilePath, IMetadataService metadataService, HttpClient httpClient, ILogger<AutoDeleteService> logger)
		{
            this.stateFilePath = stateFilePath;
            this.metadataService = metadataService;
            this.httpClient = httpClient;
            this.logger = logger;

            InitializeAutoDeleteTime();
        }

        private void InitializeAutoDeleteTime()
        {
            if (!File.Exists(stateFilePath))
            {
                AutoDeleteTime = DateTime.Now.AddHours(24);
            }
            else
            {
                // Read the date from the file.
                var fileDate = DateTime.Parse(File.ReadAllText(stateFilePath).Trim());
                // Set the ActionTime to be 24 hours from the read date.
                AutoDeleteTime = fileDate.AddHours(24);
            }
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

