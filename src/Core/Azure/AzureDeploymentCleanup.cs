using System;
using Azure;
using Azure.Core;
using Azure.ResourceManager;
using Azure.ResourceManager.Resources;
using Azure.ResourceManager.Resources.Models;
using System.Threading.Tasks;
using Microsoft.Extensions.Logging;

namespace Modm.Azure
{
    public class AzureDeploymentCleanup
	{
        private readonly IAzureResourceManagerClient azureResourceManager;
        private readonly ILogger<AzureDeploymentCleanup> logger;

        public AzureDeploymentCleanup(IAzureResourceManagerClient azureResourceManager, ILogger<AzureDeploymentCleanup> logger)
		{
            this.azureResourceManager = azureResourceManager;
            this.logger = logger;
        }

        public async Task<bool> DeleteResourcePostDeployment(string resourceGroupName)
        {
            string[] deletePhases = new string[] { "standard", "post" };

            foreach (string currentPhase in deletePhases)
            {
                this.logger.LogInformation($"deleting resources in {resourceGroupName} with tag {currentPhase}");
                bool deleted = await DeleteResourcesWithPhaseTag(resourceGroupName, currentPhase);

                if (!deleted)
                {
                    return false;
                }
            }

            return true;
        }

        private async Task<bool> DeleteResourcesWithPhaseTag(string resourceGroupName, string phase)
        {
            var resourcesToDelete = await this.azureResourceManager.GetResourcesToDeleteAsync(resourceGroupName, phase);
            int maxAttempts = resourcesToDelete.Count * 5;
            int attempt = 0;

            while (resourcesToDelete.Count > 0 && attempt < maxAttempts)
            {
                var resource = resourcesToDelete[0];

                this.logger.LogInformation($"Attempting to delete {resource}");
                if (await this.azureResourceManager.TryDeleteResourceAsync(resource))
                {
                    this.logger.LogInformation("Delete succeeded");
                    resourcesToDelete.RemoveAt(0);
                }
                else
                {
                    this.logger.LogInformation("Delete failed.  Moving to end of list");
                    resourcesToDelete.RemoveAt(0);
                    resourcesToDelete.Add(resource);
                }

                attempt++;
            }

            return (resourcesToDelete.Count == 0);
        }
    }
}

