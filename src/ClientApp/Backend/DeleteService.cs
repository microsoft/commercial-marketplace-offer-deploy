using System;
using System.Threading;
using Microsoft.Extensions.Options;
using Modm.Azure;

namespace ClientApp.Backend
{
	public class DeleteService : BackgroundService
    {
        bool deleteStarted;
        string resourceGroupName;

        private readonly AzureDeploymentCleanup cleanup;
        private readonly IConfiguration configuration;
        private ILogger<DeleteService> logger;

        private const string DeleteFileName = "delete.txt";
        private const string DeleteFileDirectoryKey = "DeleteFileDirectory";

        const int DefaultWaitDelaySeconds = 30;
        
        public DeleteService(AzureDeploymentCleanup cleanup, IConfiguration configuration, ILogger<DeleteService> logger)
		{
            this.cleanup = cleanup;
            this.configuration = configuration;
            this.logger = logger;
		}

        protected override async Task ExecuteAsync(CancellationToken cancellationToken)
        {
            this.logger.LogInformation("Waiting for delete...");
            await WaitForDelete(cancellationToken);


            if (!cancellationToken.IsCancellationRequested)
            {
                this.logger.LogInformation($"Calling DeleteResourcePostDeployment with {this.resourceGroupName}");
                await this.cleanup.DeleteResourcePostDeployment(this.resourceGroupName);
            }
        }

        async Task WaitForDelete(CancellationToken cancellationToken)
        {
            while (!deleteStarted || String.IsNullOrEmpty(this.resourceGroupName))
            {
                await Task.Delay(DefaultWaitDelaySeconds * 1000, cancellationToken);

                string stateFileContent = ReadStateFile();
                if (!String.IsNullOrEmpty(stateFileContent))
                {
                    this.logger.LogInformation("State file read");
                    this.deleteStarted = true;
                    this.resourceGroupName = stateFileContent;
                }
            }
        }

        public void Start(string resourceGroupName)
        {
            this.resourceGroupName = resourceGroupName;
            this.deleteStarted = true;

            WriteStateFile(resourceGroupName);
        }

        private string ReadStateFile()
        {
            string filePath = Path.Combine(this.configuration[DeleteFileDirectoryKey], DeleteFileName);

            if (File.Exists(filePath))
            {
                return File.ReadAllText(filePath);
            }

            return null;
        }

        private void WriteStateFile(string resourceGroupName)
        {
            string filePath = Path.Combine(this.configuration[DeleteFileDirectoryKey], DeleteFileName);
            Directory.CreateDirectory(Path.GetDirectoryName(filePath));
            File.WriteAllText(filePath, resourceGroupName);
        }
    }
}

