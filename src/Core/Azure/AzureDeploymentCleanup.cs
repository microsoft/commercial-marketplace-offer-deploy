using System;
using Azure;
using Azure.Core;
using Azure.ResourceManager;
using Azure.ResourceManager.Resources;
using Azure.ResourceManager.Resources.Models;
using System.Threading.Tasks;


namespace Modm.Azure
{
    public class AzureDeploymentCleanup
	{
        private readonly IAzureResourceManagerClient azureResourceManager;

        public AzureDeploymentCleanup(IAzureResourceManagerClient azureResourceManager)
		{
            this.azureResourceManager = azureResourceManager;
        }

        public async Task<bool> DeleteResourcePostDeployment(string resourceGroupName)
        {
            string[] deletePhases = new string[] { "standard", "post" };

            foreach (string currentPhase in deletePhases)
            {
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

                if (await this.azureResourceManager.TryDeleteResourceAsync(resource))
                {
                    resourcesToDelete.RemoveAt(0);
                }
                else
                {
                    // If deletion fails, move the resource to the end of the list
                    resourcesToDelete.RemoveAt(0);
                    resourcesToDelete.Add(resource);
                }

                attempt++;
            }

            return (resourcesToDelete.Count == 0);
        }
    }
}

