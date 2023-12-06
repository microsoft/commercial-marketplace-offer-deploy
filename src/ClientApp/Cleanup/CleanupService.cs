using System;
using Azure.Core;
using Azure.ResourceManager;
using Azure.ResourceManager.Compute;
using Azure.ResourceManager.Resources;
using MediatR;
using Modm.Azure;

namespace ClientApp.Cleanup
{
	public class CleanupService : ICleanupService
	{
        private readonly IAzureResourceManagerClient azureResourceManagerClient;
        private readonly IMediator mediator;
        private readonly ILogger<CleanupService> logger;
        private readonly string standardTag = "standard";
        private readonly string postTag = "post";

		public CleanupService(IAzureResourceManagerClient azureResourceManagerClient, IMediator mediator, ILogger<CleanupService> logger)
		{
            this.azureResourceManagerClient = azureResourceManagerClient;
            this.mediator = mediator;
            this.logger = logger;
		}

        public async Task<bool> CleanupInstallAsync(string resourceGroup)
        {
            bool success = true;

            var resourceGroupResource = await this.azureResourceManagerClient.GetResourceGroupResourceAsync(resourceGroup);

            var standardResourcesToDelete = await this.azureResourceManagerClient.GetResourcesToDeleteAsync(resourceGroup, this.standardTag);
            var standardCleanupOperations = new List<Func<Task<DeleteResourceResult>>>
            {
                () => Cleanup<DeleteVirtualMachine>("Microsoft.Compute", "virtualMachines", resourceGroupResource, standardResourcesToDelete),
                () => Cleanup<DeleteVirtualNetwork>("Microsoft.Network", "virtualNetworks", resourceGroupResource, standardResourcesToDelete),
                () => Cleanup<DeleteNetworkSecurityGroup>("Microsoft.Network", "networkSecurityGroups", resourceGroupResource, standardResourcesToDelete),
                () => Cleanup<DeleteAppConfiguration>("Microsoft.AppConfiguration", "configurationStores", resourceGroupResource, standardResourcesToDelete),
                () => Cleanup<DeleteStorageAccount>("Microsoft.Storage", "storageAccounts", resourceGroupResource, standardResourcesToDelete),
            };
            await ProcessCleanupOperations(standardCleanupOperations);

            var postResourcesToDelete = await this.azureResourceManagerClient.GetResourcesToDeleteAsync(resourceGroup, this.postTag);
            var postCleanupOperations = new List<Func<Task<DeleteResourceResult>>>
            {
                () => Cleanup<DeleteAppService>("Microsoft.Web", "sites", resourceGroupResource, postResourcesToDelete)
            };
            await ProcessCleanupOperations(postCleanupOperations);

            return success;
        }

        private async Task ProcessCleanupOperations(List<Func<Task<DeleteResourceResult>>> cleanupOperations)
        {
            foreach (var operation in cleanupOperations)
            {
                var result = await operation();
                if (!result.Succeeded)
                {
                    logger.LogError($"A cleanup operation was not successful");
                }
            }
        }

        private async Task<DeleteResourceResult> Cleanup<TCommand>(
            string resourceNamespace,
            string resourceType,
            ResourceGroupResource resourceGroupResource,
            List<GenericResource> cleanupResources) where TCommand : IDeleteResourceRequest
        {
            var resource = FindResource(resourceNamespace, resourceType, cleanupResources);
            if (resource != null)
            {
                var deleteCommand = CreateCommand<TCommand>(resource, resourceGroupResource);
                return await mediator.Send(deleteCommand);
            }

            return new DeleteResourceResult { Succeeded = true };
        }

        private GenericResource FindResource(string resourceNamespace, string resourceType, List<GenericResource> resources)
        {
            var foundResource = resources.FirstOrDefault(x => x.Data.ResourceType.Namespace.Equals(resourceNamespace)
                && x.Data.ResourceType.Type.Equals(resourceType));

            return foundResource;
        }

        private T CreateCommand<T>(
            GenericResource resource,
            ResourceGroupResource resourceGroupResource) where T : IDeleteResourceRequest
        {
            var constructorInfo = typeof(T).GetConstructor(new Type[] { typeof(ResourceGroupResource), typeof(ResourceIdentifier) });

            if (constructorInfo == null)
            {
                throw new InvalidOperationException($"The type {typeof(T).Name} does not have a constructor with the required signature.");
            }

            var command = (T)constructorInfo.Invoke(new object[] { resourceGroupResource, resource.Id });
            return command;
        }
    }
}

