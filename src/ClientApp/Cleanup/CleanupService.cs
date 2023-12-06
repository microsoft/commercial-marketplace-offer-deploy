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

            var standardResourcesToDelete = await this.azureResourceManagerClient.GetResourcesToDeleteAsync(resourceGroup, "standard");
           
            var standardCleanupOperations = new List<Func<Task<DeleteResourceResult>>>
            {
                () => CleanupVirtualMachine(resourceGroupResource, standardResourcesToDelete),
                () => CleanupVirtualNetwork(resourceGroupResource, standardResourcesToDelete),
                () => CleanupNetworkSecurityGroup(resourceGroupResource, standardResourcesToDelete),
                () => CleanupAppConfiguration(resourceGroupResource, standardResourcesToDelete),
                () => CleanupStorageAccount(resourceGroupResource, standardResourcesToDelete)
            };

            foreach (var cleanupOperation in standardCleanupOperations)
            {
                var result = await cleanupOperation();
                if (!result.Succeeded)
                {
                    //TODO: handle failure
                }
            }

            var postResourcesToDelete = await this.azureResourceManagerClient.GetResourcesToDeleteAsync(resourceGroup, "post");

            var postCleanupOperations = new List<Func<Task<DeleteResourceResult>>>
            {
                () => CleanupAppService(resourceGroupResource, standardResourcesToDelete)
            };

            foreach (var cleanupOperation in postCleanupOperations)
            {
                var result = await cleanupOperation();
                if (!result.Succeeded)
                {
                    //TODO: handle failure
                }
            }

            return success;
        }

        private async Task<DeleteResourceResult> CleanupStorageAccount(ResourceGroupResource resourceGroupResource, List<GenericResource> resourcesToDelete)
        {
            var storageAccountResource = FindResource("Microsoft.Storage", "storageAccounts", resourcesToDelete);
            var deleteStorageAccountCommand = CreateCommand<DeleteVirtualMachine>(storageAccountResource, resourceGroupResource);
            var deleteStorageAccountResult = await this.mediator.Send(deleteStorageAccountCommand);
            return deleteStorageAccountResult;
        }

        private async Task<DeleteResourceResult> CleanupAppConfiguration(ResourceGroupResource resourceGroupResource, List<GenericResource> resourcesToDelete)
        {
            var appConfigurationResource = FindResource("Microsoft.AppConfiguration", "configurationStores", resourcesToDelete);
            var deleteAppConfigurationCommand = CreateCommand<DeleteVirtualMachine>(appConfigurationResource, resourceGroupResource);
            var deleteAppConfigurationResult = await this.mediator.Send(deleteAppConfigurationCommand);
            return deleteAppConfigurationResult;
        }

        private async Task<DeleteResourceResult> CleanupAppService(ResourceGroupResource resourceGroupResource, List<GenericResource> resourcesToDelete)
        {
            var appServiceResource = FindResource("Microsoft.Web", "sites", resourcesToDelete);
            var deleteAppServiceCommand = CreateCommand<DeleteVirtualMachine>(appServiceResource, resourceGroupResource);
            var deleteAppServiceResult = await this.mediator.Send(deleteAppServiceCommand);
            return deleteAppServiceResult;
        }

        private async Task<DeleteResourceResult> CleanupNetworkSecurityGroup(ResourceGroupResource resourceGroupResource, List<GenericResource> resourcesToDelete)
        {
            var nsgResource = FindResource("Microsoft.Network", "networkSecurityGroups", resourcesToDelete);
            var deleteNsgCommand = CreateCommand<DeleteVirtualMachine>(nsgResource, resourceGroupResource);
            var deleteNsgResult = await this.mediator.Send(deleteNsgCommand);
            return deleteNsgResult;
        }

        private async Task<DeleteResourceResult> CleanupVirtualMachine(ResourceGroupResource resourceGroupResource, List<GenericResource> resourcesToDelete)
        {
            var virtualMachineResource = FindResource("Microsoft.Compute", "virtualMachines", resourcesToDelete);
            var deleteVmCommand = CreateCommand<DeleteVirtualMachine>(virtualMachineResource, resourceGroupResource);
            var deleteVmResult = await this.mediator.Send(deleteVmCommand);
            return deleteVmResult;
        }

        private async Task<DeleteResourceResult> CleanupVirtualNetwork(ResourceGroupResource resourceGroupResource, List<GenericResource> resourcesToDelete)
        {
            var virtualNetworkResource = FindResource("Microsoft.Network", "virtualNetworks", resourcesToDelete);
            var deleteVnetCommand = CreateCommand<DeleteVirtualNetwork>(virtualNetworkResource, resourceGroupResource);
            var deleteVnetResult = await this.mediator.Send(deleteVnetCommand);
            return deleteVnetResult;
        }

        private GenericResource FindResource(string resourceNamespace, string resourceType, List<GenericResource> resources)
        {
            var foundResource = resources.FirstOrDefault(x => x.Data.ResourceType.Namespace.Equals(resourceNamespace)
                && x.Data.ResourceType.Type.Equals(resourceType));

            return foundResource;
        }

        private T CreateCommand<T>(GenericResource resource, ResourceGroupResource resourceGroupResource) where T : IDeleteResourceRequest
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

