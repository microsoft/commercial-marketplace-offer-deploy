using System;
using Azure;
using Azure.ResourceManager;
using Azure.ResourceManager.Resources;

namespace Modm.Azure
{
    public class AzureResourceManagerClient : IAzureResourceManagerClient
    {
        private readonly ArmClient client;

        public AzureResourceManagerClient(ArmClient client)
        {
            this.client = client;
        }

        public async Task<List<GenericResource>> GetResourcesToDeleteAsync(string resourceGroupName, string phase)
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

        public async Task<ResourceGroupResource> GetResourceGroupResourceAsync(string resourceGroupName)
        {
            var subscription = await client.GetDefaultSubscriptionAsync();
            var response = await subscription.GetResourceGroupAsync(resourceGroupName);
            var resourceGroup = response.Value;
            return resourceGroup;
        }

        public async Task<bool> TryDeleteResourceAsync(GenericResource resource)
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

