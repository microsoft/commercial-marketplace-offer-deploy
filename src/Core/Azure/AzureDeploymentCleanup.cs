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
		private readonly ArmClient client;

        public AzureDeploymentCleanup(ArmClient client)
		{
			this.client = client;
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
            var resourcesToDelete = await GetResourcesToDelete(resourceGroupName, phase);
            int maxAttempts = resourcesToDelete.Count * 5;
            int attempt = 0;

            while (resourcesToDelete.Count > 0 && attempt < maxAttempts)
            {
                var resource = resourcesToDelete[0];

                if (await TryDeleteResource(resource))
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

        private async Task<List<GenericResource>> GetResourcesToDelete(string resourceGroupName, string phase)
        {
            var subscription = await client.GetDefaultSubscriptionAsync();
            var response = await subscription.GetResourceGroupAsync(resourceGroupName);
            var resourceGroup = response.Value;

            var resourcesToDelete = new List<GenericResource>();

            await foreach (var resource in resourceGroup.GetGenericResourcesAsync())
            {
                if (resource.Data.Tags != null && resource.Data.Tags.TryGetValue("modm", out var tagValue) && tagValue == phase)
                {
                    resourcesToDelete.Add(resource);
                }
            }

            return resourcesToDelete;
        }

        private async Task<bool> TryDeleteResource(GenericResource resource)
        {
            try
            {
                await resource.DeleteAsync(WaitUntil.Started);
                return true;
            }
            catch
            {
                return false; // Return false if deletion fails
            }
        }
    }
}

