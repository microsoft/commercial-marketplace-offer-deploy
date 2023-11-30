using System;
using System.Threading;
using Modm.Azure;

namespace ClientApp.Backend
{
	public class DeleteService : BackgroundService
    {
        bool controllerStarted;
        string resourceGroupName;

        private readonly AzureDeploymentCleanup cleanup;
        private const string DataDirectory = "/home/site/wwwroot/data";
        private const string DeleteFileName = "delete.txt";
        const int DefaultWaitDelaySeconds = 30;

        public DeleteService(IAzureResourceManager resourceManager)
		{
            this.cleanup = new AzureDeploymentCleanup(resourceManager);
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
            while (!controllerStarted)
            {
                await Task.Delay(DefaultWaitDelaySeconds * 1000, cancellationToken);

                if (DeleteFileExists())
                {
                    this.controllerStarted = true;
                }
            }
        }

        public void Start()
        {
            this.controllerStarted = true;
            WriteStateFile();
        }

        private bool DeleteFileExists()
        {
            string filePath = Path.Combine(DataDirectory, DeleteFileName);
            return File.Exists(filePath);
        }

        private void WriteStateFile()
        {
            string content = $"Delete initiated - {DateTime.UtcNow:O}";
            string filePath = Path.Combine(DataDirectory, DeleteFileName);
            Directory.CreateDirectory(Path.GetDirectoryName(filePath));
            File.WriteAllText(filePath, content);
        }
    }
}

