using System;
using Azure.ResourceManager.Resources;

namespace Modm.Azure
{
    public interface IAzureResourceManagerClient
    {
        Task<List<GenericResource>> GetResourcesToDeleteAsync(string resourceGroupName, string phase);
        Task<bool> TryDeleteResourceAsync(GenericResource resource);
    }
}

