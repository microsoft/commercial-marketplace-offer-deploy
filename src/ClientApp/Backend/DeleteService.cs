using System;
using System.Threading;
using Microsoft.Extensions.Options;
using Modm.Azure;

namespace ClientApp.Backend
{
	public class DeleteService : BackgroundService
    {
        bool controllerStarted;
        string resourceGroupName;

        private readonly DeleteServiceOptions options;
        private readonly AzureDeploymentCleanup cleanup;

        private const string DeleteFileName = "delete.txt";

        const int DefaultWaitDelaySeconds = 30;
        
        public DeleteService(IAzureResourceManager resourceManager, IOptions<DeleteServiceOptions> options)
		{
            this.cleanup = new AzureDeploymentCleanup(resourceManager);
            this.options = options.Value;
		}

        protected override async Task ExecuteAsync(CancellationToken cancellationToken)
        {
            await WaitForControllerToStart(cancellationToken);

            if (!cancellationToken.IsCancellationRequested)
            {
                await this.cleanup.DeleteResourcePostDeployment(this.resourceGroupName);
            }
        }

        async Task WaitForControllerToStart(CancellationToken cancellationToken)
        {
            while (!controllerStarted || String.IsNullOrEmpty(this.resourceGroupName))
            {
                await Task.Delay(DefaultWaitDelaySeconds * 1000, cancellationToken);

                string stateFileContent = ReadStateFile();
                if (!String.IsNullOrEmpty(stateFileContent))
                {
                    this.controllerStarted = true;
                    this.resourceGroupName = stateFileContent;
                }
            }
        }

        public void Start(string resourceGroupName)
        {
            this.resourceGroupName = resourceGroupName;
            this.controllerStarted = true;
            WriteStateFile(resourceGroupName);
        }

        private string ReadStateFile()
        {
            string filePath = Path.Combine(this.options.DataDirectory, DeleteFileName);

            if (File.Exists(filePath))
            {
                return File.ReadAllText(filePath);
            }

            return null;
        }

        private void WriteStateFile(string resourceGroupName)
        {
            string filePath = Path.Combine(this.options.DataDirectory, DeleteFileName);
            Directory.CreateDirectory(Path.GetDirectoryName(filePath));
            File.WriteAllText(filePath, resourceGroupName);
        }
    }
}

