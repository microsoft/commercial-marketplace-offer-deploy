using System;
using System.Threading;
using Azure.Core;
using Azure.ResourceManager;
using Azure.ResourceManager.Compute;
using Azure.ResourceManager.Resources;
using MediatR;
using Modm.Azure;

namespace ClientApp.Cleanup
{
	public class DeleteProcessor : IDeleteProcessor
	{
        public const string StandardTag = "standard";
        public const string PostTag = "post";

        private readonly ArmClient armClient;
        private readonly IMediator mediator;
        private readonly ILogger<DeleteProcessor> logger;
        
		public DeleteProcessor(ArmClient azureResourceManagerClient, IMediator mediator, ILogger<DeleteProcessor> logger)
		{
            this.armClient = azureResourceManagerClient;
            this.mediator = mediator;
            this.logger = logger;
		}

        public async Task<bool> DeleteInstallResourcesAsync(string resourceGroup, CancellationToken cancellationToken)
        {
            bool success = true;

            var resourceGroupResource = await this.GetResourceGroupResourceAsync(resourceGroup);

            var standardResourcesToDelete = await this.GetResourcesToDeleteAsync(resourceGroup, DeleteProcessor.StandardTag);
            var standardCleanupOperations = new List<Func<Task<DeleteResourceResult>>>
            {
                () => Cleanup<DeleteVirtualMachine>("Microsoft.Compute", "virtualMachines", resourceGroupResource, standardResourcesToDelete, cancellationToken),
                () => Cleanup<DeleteVirtualNetwork>("Microsoft.Network", "virtualNetworks", resourceGroupResource, standardResourcesToDelete, cancellationToken),
                () => Cleanup<DeleteNetworkSecurityGroup>("Microsoft.Network", "networkSecurityGroups", resourceGroupResource, standardResourcesToDelete, cancellationToken),
                () => Cleanup<DeleteAppConfiguration>("Microsoft.AppConfiguration", "configurationStores", resourceGroupResource, standardResourcesToDelete, cancellationToken),
                () => Cleanup<DeleteStorageAccount>("Microsoft.Storage", "storageAccounts", resourceGroupResource, standardResourcesToDelete, cancellationToken),
            };
            await ProcessCleanupOperations(standardCleanupOperations);

            var postResourcesToDelete = await this.GetResourcesToDeleteAsync(resourceGroup, DeleteProcessor.PostTag);
            var postCleanupOperations = new List<Func<Task<DeleteResourceResult>>>
            {
                () => Cleanup<DeleteAppService>("Microsoft.Web", "sites", resourceGroupResource, postResourcesToDelete, cancellationToken)
            };
            await ProcessCleanupOperations(postCleanupOperations);

            return success;
        }

        protected virtual async Task<List<GenericResource>> GetResourcesToDeleteAsync(string resourceGroupName, string phase)
        {
            var subscription = await this.armClient.GetDefaultSubscriptionAsync();
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

        protected virtual async Task<ResourceGroupResource> GetResourceGroupResourceAsync(string resourceGroupName)
        {
            var subscription = await this.armClient.GetDefaultSubscriptionAsync();
            var response = await subscription.GetResourceGroupAsync(resourceGroupName);
            var resourceGroup = response.Value;
            return resourceGroup;
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
            List<GenericResource> cleanupResources,
            CancellationToken cancellationToken) where TCommand : IDeleteResourceRequest
        {
            var resource = FindResource(resourceNamespace, resourceType, cleanupResources);
            if (resource != null)
            {
                var deleteCommand = CreateCommand<TCommand>(resource, resourceGroupResource);
                return await mediator.Send(deleteCommand, cancellationToken);
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

